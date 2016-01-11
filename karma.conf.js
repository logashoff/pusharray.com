module.exports = function(config) {
  config.set({
    // base path, that will be used to resolve files and exclude
    basePath: './',

    preprocessors: {
      'static/src/**/*.ts': 'browserify',
      'static/test/**/*.ts': 'browserify'
    },

    frameworks: ['browserify', 'jasmine'],

    browserify: {
      debug: true,
      plugin: [
        ['tsify', {target: 'es5'}]
      ]
    },

    files: ['static/test/**/*_spec.ts'],

    reporters: ['progress'],

    // web server port
    port: 9876,

    // enable / disable colors in the output (reporters and logs)
    colors: true,

    // level of logging
    logLevel: config.LOG_INFO,

    // enable / disable watching file and executing tests whenever any file changes
    autoWatch: false,

    // Start these browsers
    browsers: ['PhantomJS'],

    // If browser does not capture in given timeout [ms], kill it
    captureTimeout: 60000,

    // Continuous Integration mode
    // if true, it capture browsers, run tests and exit
    singleRun: true
  });
};
