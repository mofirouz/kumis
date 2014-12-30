(function() {
  Polymer('kumis-topics', {
    publish: {
      zkClusterAddress: "" // init blank
    },
    computed: {
      zk: 'computeKumisBrokerAddress(zkClusterAddress)'
    },
    computeKumisBrokerAddress: function(str) {
      var hostname = window.location.host.substring(0, window.location.host.lastIndexOf(":"));
      return window.location.protocol + "//" + hostname + ":7777"  + "/" + str;
    }
  });
})();