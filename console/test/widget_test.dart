import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/main.dart';
import 'package:retrovibed/retrovibed.dart' as retro;

void main() {
  testWidgets('smoke test', (WidgetTester tester) async {
    // bind the embedded daemon to an OS-assigned port and isolate its
    // storage under a scratch dir, so this test never collides with a
    // dev instance already running on the machine (fixed :9998, shared
    // meta.db lock).
    final xdg = Directory.systemTemp.createTempSync('retrovibed_widget_test_');
    addTearDown(() => xdg.deleteSync(recursive: true));
    retro.setenv('RETROVIBED_DAEMON_SOCKET', 'tcp://:0');
    retro.setenv('XDG_CONFIG_HOME', xdg.path);
    retro.setenv('XDG_CACHE_HOME', xdg.path);
    retro.setenv('XDG_DATA_HOME', xdg.path);

    // Build our app and trigger a frame.
    await tester.pumpWidget(Retrovibed());
    expect(tester.takeException(), isNull);
  });
}
