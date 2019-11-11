
export default {
  mode: 'universal',
  /*
  ** Headers of the page
  */
  head: {
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: process.env.npm_package_description || '' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      {
        rel: "stylesheet",
        href:
          "https://fonts.googleapis.com/css?family=Permanent+Marker&display=swap"
      }
    ]
  },
  /*
  ** Customize the progress-bar color
  */
  loading: { color: '#fff' },
  /*
  ** Global CSS
  */
  css: [
  ],
  /*
  ** Plugins to load before mounting the App
  */
  plugins: [
    '~/plugins/projects.server.js',
    '~/plugins/htmlDecode.js'
  ],
  /*
  ** Nuxt.js dev-modules
  */
  buildModules: [
  ],
  /*
  ** Nuxt.js modules
  */
  modules: [
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/axios',
    '@nuxtjs/proxy',
    '@nuxtjs/date-fns',
    'nuxt-buefy',
    ['nuxt-vuex-localstorage', {
      localStorage: ['jobs'],
      sessionStorage: ['jobs']
    }],
    ['nuxt-env', {
      keys: [
        'REQUEST_TOKEN',
        // { key: 'REQUEST_TOKEN', secret: true } // Only inject the var server side
        'BUILD_HASH'
      ]
    }]
    // '@nuxtjs/style-resources'
  ],
  /*
  ** Axios module configuration
  ** See https://axios.nuxtjs.org/options
  */
  axios: {
    prefix: '/api',
    proxy: true
  },
  proxy: {
    '/api/': {target: process.env.BACKEND_URL || 'http://localhost:8000/', pathRewrite: {'^/api/': ''}}
  },
  // styleResources: {
  //   // your settings here
  //   scss: [
  //    '@/assets/styles/_variables.scss',
  //    '@/assets/styles/_mixins.scss',
  //   ],
  //  },
  /*
  ** Build configuration
  */
  build: {
    /*
    ** You can extend webpack config here
    */
   transpile: ['vuejs-smart-table'],
    // postcss: {
    //   preset: {
    //     features: {
    //       customProperties: false
    //     }
    //   }
    // },
    extend (config, ctx) {
    }
  }
}
