function(
  ubuntu_version='1604-xenial',
  additional_licenses=[],
) {
  Steps: {
    create_image: {
      CreateImages: [
        {
          Licenses: std.setUnion(
            ['projects/ubuntu-os-cloud/global/licenses/ubuntu-%s' % ubuntu_version],
            additional_licenses,
          ),
        },
      ],
    },
  },
}
