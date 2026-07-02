import './api.dart' as api;

final Storage = api.Quota(
  sku: "59790d9e-9c74-152c-fceb-a26607f02146",
  description: "Storage",
  adjustable: true,
);

final Bandwidth = api.Quota(
  sku: "620cac69-da27-2e19-ae99-d162275e528d",
  description: "Bandwidth",
  adjustable: true,
);

final Profiles = api.Quota(
  sku: "4a76f4f7-a114-5b16-81f8-cd2c19958be5",
  description: "allow number of users",
  adjustable: false,
);
