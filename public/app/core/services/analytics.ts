import $ from 'jquery';
import coreModule from 'app/core/core_module';
import config from 'app/core/config';

export class Analytics {
  /** @ngInject */
  constructor(private $rootScope, private $location) {}

  gaInit() {
    $.ajax({
      url: 'https://www.google-analytics.com/analytics.js',
      dataType: 'script',
      cache: true,
    });
    const ga = ((window as any).ga =
      (window as any).ga ||
      //tslint:disable-next-line:only-arrow-functions
      function() {
        (ga.q = ga.q || []).push(arguments);
      });
    ga.l = +new Date();
    ga('create', (config as any).googleAnalyticsId, 'auto');
    ga('set', 'anonymizeIp', true);
    return ga;
  }

  init() {
    this.$rootScope.$on('$viewContentLoaded', () => {
      const track = { page: this.$location.url() };
      const ga = (window as any).ga || this.gaInit();
      ga('set', track);

      const userEmail = this.getUserEmail();
      if (userEmail) {
        ga('set', 'userId', userEmail); // Set the user ID using signed-in user_id.
      }

      ga('send', 'pageview');
    });
  }

  getUserEmail() {
    var bootData = (<any>window).grafanaBootData || { settings: {} };
    if (bootData.user) {
      return bootData.user.email;
    }

    return null;
  }
}

/** @ngInject */
function startAnalytics(googleAnalyticsSrv) {
  if ((config as any).googleAnalyticsId) {
    googleAnalyticsSrv.init();
  }
}

coreModule.service('googleAnalyticsSrv', Analytics).run(startAnalytics);
