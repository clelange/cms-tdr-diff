package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

// Configuration structure
type Configuration struct {
	port                  int
	frontendOrigin        string
	gitlabToken           string
	triggerToken          string
	gitlabURL             string
	gitlabProject         int
	debug                 bool
	prefetchGroups        bool
	groupIds              []string
	commitHistoryDays     int
	updateIntervalSeconds int
}

type tdrTypes struct {
	Names []string `json:"names"`
}

type group struct {
	ID string `uri:"id" binding:"required"`
}

type projectStruct struct {
	Group string `uri:"group" binding:"required"`
	ID    string `uri:"id" binding:"required"`
}

type triggerStruct struct {
	Project string `uri:"project" binding:"required"`
	Group   string `uri:"group" binding:"required"`
	SHA1    string `uri:"sha1" binding:"required"`
	SHA2    string `uri:"sha2" binding:"required"`
}

type pipelineStruct struct {
	ID int `uri:"id" binding:"required"`
}

type gitlabProjectList struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	WebURL         string     `json:"web_url"`
	LastActivityAt *time.Time `json:"last_activity_at"`
	Description    string     `json:"description"`
}

type gitlabCommitList struct {
	ID          string     `json:"id"`
	ShortID     string     `json:"short_id"`
	CreatedAt   *time.Time `json:"created_at"`
	Title       string     `json:"title"`
	AuthorName  string     `json:"author_name"`
	AuthorEmail string     `json:"author_email"`
}

// type gitlabPipelineStatus struct {
// 	ID         int    `json:"id"`
// 	Status     string `json:"status"`
// 	Ref        string `json:"ref"`
// 	SHA        string `json:"sha"`
// 	BeforeSHA  string `json:"before_sha"`
// 	Tag        bool   `json:"tag"`
// 	YamlErrors string `json:"yaml_errors"`
// 	User       struct {
// 		Name      string `json:"name"`
// 		Username  string `json:"username"`
// 		ID        int    `json:"id"`
// 		State     string `json:"state"`
// 		AvatarURL string `json:"avatar_url"`
// 		WebURL    string `json:"web_url"`
// 	}
// 	UpdatedAt   *time.Time `json:"updated_at"`
// 	CreatedAt   *time.Time `json:"created_at"`
// 	StartedAt   *time.Time `json:"started_at"`
// 	FinishedAt  *time.Time `json:"finished_at"`
// 	CommittedAt *time.Time `json:"committed_at"`
// 	Duration    int        `json:"duration"`
// 	Coverage    string     `json:"coverage"`
// 	WebURL      string     `json:"web_url"`
// }

func readConfig() (*viper.Viper, error) {
	v := viper.New()
	// Use envrironment variables for hosted version with VIPER_ prefix
	v.SetEnvPrefix("viper") // will be uppercased automatically
	// v.BindEnv("gitlabToken")
	// v.BindEnv("triggerToken")
	v.SetDefault("port", 8000)
	v.SetDefault("frontendOrigin", "http://localhost:3000")
	v.SetDefault("gitlabURL", "https://gitlab.cern.ch/api/v4")
	v.SetDefault("gitlabProject", 56283)
	v.SetDefault("debug", true)
	v.SetDefault("prefetchGroups", true)
	v.SetDefault("commitHistoryDays", 90)
	v.SetDefault("updateIntervalSeconds", 300)
	v.SetDefault("groupIds", []string{
		"papers", "notes", "reports",
	})
	v.SetConfigName("config")
	v.AddConfigPath("config")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config file not found. Trying to use environment variables.")
			return v, nil
		}
	}
	return v, err
}

// check that provided subgroups exist in project
func validateSubgroups(groupID int, gl *gitlab.Client, configuration Configuration) (map[string]int, error) {
	groups, _, err := gl.Groups.ListSubgroups(groupID, nil)
	if err != nil {
		fmt.Print(err)
	}

	groupIDMap := make(map[string]int)
	for _, n := range configuration.groupIds {
		found := false
		for _, group := range groups {
			if n == group.Name {
				groupIDMap[group.Name] = group.ID
				found = true
			}
		}
		if !found {
			err = errors.New("subgroup name not found")
			return groupIDMap, err
		}
	}
	return groupIDMap, err
}

// get all projects for a given subgroup
func getProjects(groupID int, gl *gitlab.Client) ([]gitlabProjectList, error) {
	maxPages := 100
	currentPage := 1
	var projectList []gitlabProjectList
	for currentPage <= maxPages {
		var listQueryOptions = &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100, // this is the maximum one can ask for
				Page:    currentPage,
			}}
		projects, response, err := gl.Groups.ListGroupProjects(groupID, listQueryOptions)
		if err != nil {
			fmt.Print(err)
			return projectList, err
		}
		currentProjectList := make([]gitlabProjectList, len(projects))
		for i := 0; i < len(projects); i++ {
			currentProjectList[i] = gitlabProjectList{
				ID:             projects[i].ID,
				Name:           projects[i].Name,
				WebURL:         projects[i].WebURL,
				LastActivityAt: projects[i].LastActivityAt,
				Description:    projects[i].Description,
			}
			projectList = append(projectList, currentProjectList[i])
		}
		maxPages = response.TotalPages
		currentPage++
	}
	fmt.Println("Number of projects:", len(projectList))
	return projectList, nil
}

func updateProjects(groupIDs map[string]int, gl *gitlab.Client) (map[string][]gitlabProjectList, error) {
	allProjects := make(map[string][]gitlabProjectList)
	var err error
	for value, key := range groupIDs {
		fmt.Println("Getting projects for group:", value, key)
		allProjects[value], err = getProjects(key, gl)
		if err != nil {
			fmt.Print(err)
		}
	}
	return allProjects, err
}

func getCommits(projectID int, gl *gitlab.Client) ([]gitlabCommitList, error) {
	maxPages := 100
	currentPage := 1
	var commitList []gitlabCommitList
	for currentPage <= maxPages {
		var listQueryOptions = &gitlab.ListCommitsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100, // this is the maximum one can ask for
				Page:    currentPage,
			}}
		commits, response, err := gl.Commits.ListCommits(projectID, listQueryOptions)
		if err != nil {
			fmt.Print(err)
			return commitList, err
		}
		currentCommitList := make([]gitlabCommitList, len(commits))
		for i := 0; i < len(commits); i++ {
			currentCommitList[i] = gitlabCommitList{
				ID:          commits[i].ID,
				ShortID:     commits[i].ShortID,
				CreatedAt:   commits[i].CreatedAt,
				Title:       commits[i].Title,
				AuthorName:  commits[i].AuthorName,
				AuthorEmail: commits[i].AuthorEmail,
			}
			commitList = append(commitList, currentCommitList[i])
		}
		maxPages = response.TotalPages
		currentPage++
	}
	fmt.Println("Number of commits:", len(commitList))
	return commitList, nil
}

func ping(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func main() {

	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	v1, err := readConfig()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	// fmt.Println(v1)

	var configuration Configuration
	configuration.port = v1.GetInt("port")
	configuration.frontendOrigin = v1.GetString("frontendOrigin")
	configuration.gitlabURL = v1.GetString("gitlabURL")
	configuration.gitlabProject = v1.GetInt("gitlabProject")
	configuration.debug = v1.GetBool("debug")
	configuration.prefetchGroups = v1.GetBool("prefetchGroups")
	configuration.groupIds = v1.GetStringSlice("groupIds")
	configuration.commitHistoryDays = v1.GetInt("commitHistoryDays")
	configuration.updateIntervalSeconds = v1.GetInt("updateIntervalSeconds")
	configuration.gitlabToken = v1.GetString("gitlabToken")
	if configuration.gitlabToken == "" {
		log.Panicln("gitlabToken cannot be empty.")
	}
	configuration.triggerToken = v1.GetString("triggerToken")
	if configuration.triggerToken == "" {
		log.Panicln("triggerToken cannot be empty.")
	}

	fmt.Printf("Reading config for port = %d\n", configuration.port)
	fmt.Printf("Reading config for frontendOrigin = %s\n", configuration.frontendOrigin)
	fmt.Printf("Reading config for gitlabURL = %s\n", configuration.gitlabURL)
	fmt.Printf("Reading config for gitlabProject = %d\n", configuration.gitlabProject)
	fmt.Printf("Reading config for debug = %t\n", configuration.debug)
	fmt.Printf("Reading config for prefetchGroups = %t\n", configuration.prefetchGroups)
	fmt.Printf("Reading config for groupIds = %#v\n", configuration.groupIds)
	fmt.Printf("Reading config for commitHistoryDays = %d\n", configuration.commitHistoryDays)
	fmt.Printf("Reading config for updateIntervalSeconds = %d\n", configuration.updateIntervalSeconds)
	// fmt.Printf("Reading config for gitlabToken = %s\n", configuration.gitlabToken)
	// fmt.Printf("Reading config for triggerToken = %s\n", configuration.triggerToken)

	gl := gitlab.NewClient(nil, configuration.gitlabToken)
	gl.SetBaseURL(configuration.gitlabURL)

	pipelineProject, _, err := gl.Projects.GetProject("clange/tdr-diff", nil)
	if err != nil {
		fmt.Print(err)
	}
	pipelineProjectID := pipelineProject.ID
	fmt.Println("Pipeline project ID:", pipelineProjectID)

	groupIDs, err := validateSubgroups(16284, gl, configuration) // this is the tdr group
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(groupIDs)

	allProjects, err := updateProjects(groupIDs, gl)
	// fmt.Println(allProjects)
	lastUpdated := time.Now()
	ticker := time.NewTicker(time.Duration(configuration.updateIntervalSeconds) * time.Second)
	go func() {
		for range ticker.C {
			fmt.Println("updating...")
			tempAllProjects, err := updateProjects(groupIDs, gl)
			if err != nil {
				fmt.Print("Updating projects failed", err)
				return
			}
			lastUpdated = time.Now()
			allProjects = tempAllProjects
		}
	}()

	types := make([]string, 0, len(configuration.groupIds))
	for _, key := range configuration.groupIds {
		types = append(types, key)
	}
	fmt.Println(types)
	projTypes := &tdrTypes{
		Names: types}
	resTypes, _ := json.Marshal(projTypes)
	fmt.Println(string(resTypes))
	if !configuration.debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/ping", ping)

	r.GET("/lastUpdated", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{"lastUpdated": lastUpdated})
	})

	r.GET("/types", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{"names": types})
	})

	r.GET("/projects/:id", func(c *gin.Context) {
		var groupID group
		if err := c.ShouldBindUri(&groupID); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		fmt.Println(groupID)
		for key := range groupIDs {
			fmt.Println(key)
			if key == groupID.ID {
				c.Header("Content-Type", "application/json")
				c.JSON(200, gin.H{"data": allProjects[groupID.ID]})
				return
			}
		}
		errorMessage := "Project not found: " + groupID.ID
		err = errors.New(errorMessage)
		c.JSON(404, gin.H{"msg": errorMessage})
	})

	r.GET("/commits/:group/:id", func(c *gin.Context) {
		var projectID projectStruct
		if err := c.ShouldBindUri(&projectID); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		fmt.Println(projectID)
		// get the project ID, then its commits
		projectPath := "tdr/" + projectID.Group + "/" + projectID.ID
		project, response, err := gl.Projects.GetProject(projectPath, nil)
		if err != nil {
			c.JSON(404, gin.H{"msg": response})
		}
		projectInfo := gitlabProjectList{
			ID:             project.ID,
			Name:           project.Name,
			Description:    project.Description,
			WebURL:         project.WebURL,
			LastActivityAt: project.LastActivityAt,
		}
		commitList, err := getCommits(project.ID, gl)
		if err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{"project_info": projectInfo, "commits": commitList})
	})

	r.POST("/trigger", func(c *gin.Context) {
		var triggerObject triggerStruct
		if err := c.BindJSON(&triggerObject); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		var variables = make(map[string]string)
		variables["REPO_PROJECT"] = triggerObject.Project
		variables["REPO_GROUP"] = triggerObject.Group
		variables["GIT_SHA1"] = triggerObject.SHA1
		variables["GIT_SHA2"] = triggerObject.SHA2

		referenceBranch := "master"
		pipelineOptions := &gitlab.RunPipelineTriggerOptions{
			Ref:       &referenceBranch,
			Token:     &configuration.triggerToken,
			Variables: variables,
		}

		pipeline, _, err := gl.PipelineTriggers.RunPipelineTrigger(pipelineProjectID, pipelineOptions)
		if err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{"status": "Pipeline triggered successfully!",
			"pipeline_id": pipeline.ID})
	})

	r.GET("/status/pipeline/:id", func(c *gin.Context) {
		var queryPipeline pipelineStruct
		if err := c.ShouldBindUri(&queryPipeline); err != nil {
			fmt.Println()
			c.JSON(400, gin.H{"msg": err})
			return
		}
		pipelineJobs, _, err := gl.Jobs.ListPipelineJobs(pipelineProjectID, queryPipeline.ID, nil)
		if err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		jobID := pipelineJobs[0].ID
		job, _, err := gl.Jobs.GetJob(pipelineProjectID, jobID, nil)
		if err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{"job_status": job})
	})

	// TODO: implement callback from GitLab for status update
	// TODO: Get only commits of last N days

	r.Run(":" + strconv.Itoa(configuration.port)) // listen and serve on 0.0.0.0:8000
}
