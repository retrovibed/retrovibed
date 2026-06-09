import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/retrovibed.dart' as retro;

void main() async {
  // smoke test is to ensure the daemon and the initializing code can be executed.
  // used by CI/CD as a case to limit build / environment issues across platforms.

  await retro.run(() {
    retro.seed(uuidx.random());
    retro.daemon();
  });
}
