import $ from 'jquery';
import coreModule from 'app/core/core_module';
import config from 'app/core/config';

export class Analytics {
  /** @ngInject */
  constructor(private $rootScope, private $location) {}

  gaInit() {
    $.getScript('https://www.google-analytics.com/analytics.js'); // jQuery shortcut
    var ga = ((<any>window).ga =
      (<any>window).ga ||
      function() {
        (ga.q = ga.q || []).push(arguments);
      });
    ga.l = +new Date();
    ga('create', (<any>config).googleAnalyticsId, 'auto');
    return ga;
  }

  init() {
    this.$rootScope.$on('$viewContentLoaded', () => {
      var track = { page: this.$location.url() };
      var ga = (<any>window).ga || this.gaInit();
      ga('set', track);

      const userEmail = this.getUserEmail()
      if (userEmail) {
        ga('set', 'userId', userEmail); // Set the user ID using signed-in user_id.
      }

      ga('send', 'pageview');
    });
  }

  getUserEmail() {
    var bootData = (<any>window).grafanaBootData || { settings: {} };
    if (bootData.user)
    {
      return bootData.user.email
    }

    return null
  }
}

/** @ngInject */
function startAnalytics(googleAnalyticsSrv) {
  if ((<any>config).googleAnalyticsId) {
    googleAnalyticsSrv.init();
  }
}

coreModule.service('googleAnalyticsSrv', Analytics).run(startAnalytics);
