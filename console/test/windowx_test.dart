import 'package:flutter_test/flutter_test.dart';
import 'package:window_manager/window_manager.dart' show TitleBarStyle;
import 'package:retrovibed/windowx.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  // window_manager 0.7 is built on nativeapi (ffi), there is no method channel
  // to intercept. Tests run in a VM without the native library, where the real
  // defaults must degrade to harmless no-ops with settled state.
  group('WindowManagerX real defaults without a native window', () {
    test('settle completes', () async {
      await WindowManagerX().settle();
    });

    test('addListener/removeListener complete', () async {
      final wm = WindowManagerX();
      const events = WindowEventsX();
      wm.addListener(events);
      wm.removeListener(events);
    });

    test('setTitleBarStyle completes', () async {
      await WindowManagerX().setTitleBarStyle(TitleBarStyle.hidden);
    });

    test('setFullScreen completes', () async {
      await WindowManagerX().setFullScreen(true);
    });

    test('isFullScreen reports false', () async {
      expect(await WindowManagerX().isFullScreen(), isFalse);
    });

    test('maximize/unmaximize/isMaximized complete and report false', () async {
      final wm = WindowManagerX();
      await wm.maximize();
      await wm.unmaximize();
      expect(await wm.isMaximized(), isFalse);
    });

    test('settle completes and reports a settled state', () async {
      final wm = WindowManagerX();
      await wm.settle();

      expect(await wm.isFullScreen(), isFalse);
      expect(await wm.isMaximized(), isFalse);
    });
  });

  group('WindowManagerX.fake', () {
    test('defaults to false/false without touching a channel', () async {
      final wm = WindowManagerX.fake();
      expect(await wm.isFullScreen(), isFalse);
      expect(await wm.isMaximized(), isFalse);
    });

    test('reports the configured fullscreen/maximized values', () async {
      final wm = WindowManagerX.fake(fullscreen: true, maximized: true);
      expect(await wm.isFullScreen(), isTrue);
      expect(await wm.isMaximized(), isTrue);
    });

    test('records every invoked method in order', () async {
      final calls = <String>[];
      final wm = WindowManagerX.fake(calls: calls);

      await wm.settle();
      await wm.setTitleBarStyle(TitleBarStyle.hidden);
      await wm.setFullScreen(true);
      await wm.isFullScreen();
      await wm.maximize();
      await wm.unmaximize();
      await wm.isMaximized();
      await wm.close();
      wm.addListener(const WindowEventsX());
      wm.removeListener(const WindowEventsX());

      expect(calls, [
        'settle',
        'setTitleBarStyle',
        'setFullScreen',
        'isFullScreen',
        'maximize',
        'unmaximize',
        'isMaximized',
        'close',
        'addListener',
        'removeListener',
      ]);
    });
  });
}
