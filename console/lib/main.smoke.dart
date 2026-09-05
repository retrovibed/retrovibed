import 'dart:io';

import 'package:retrovibed/retrovibed.dart' as retro;

void main() async {
  // smoke test is to ensure the daemon and the initializing code can be executed.
  // used by CI/CD as a case to limit build / environment issues across platforms.

  await retro.run(() {
    retro.setenv("RETROVIBED_SMOKE", "true");
    retro.guest();
    retro.daemon(smoke: true);
  });
  exit(0);
}
