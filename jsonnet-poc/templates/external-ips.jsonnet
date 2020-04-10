function(
  use_external_ip,
  import_network='global/networks/default',
  import_subnet='',
) {
  Steps: {
    import_disk: {
      CreateInstances: [
        {
          networkInterfaces: if use_external_ip then [
            {
              network: import_network,
              subnetwork: import_subnet,
            },
          ] else [],
        },
      ],
    },
  },
}
