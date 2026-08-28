import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:window_manager/window_manager.dart';
import 'package:retrovibed/windowx.dart';

List<String> _calls = [];
bool _fullscreen = false;
bool _maximized = false;
bool _minimized = false;

void _mockWindowManager() {
  _calls = [];
  _fullscreen = false;
  _maximized = false;
  _minimized = false;
  TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger.setMockMethodCallHandler(
    const MethodChannel('window_manager'),
    (call) async {
      _calls.add(call.method);
      switch (call.method) {
        case 'isFullScreen':
          return _fullscreen;
        case 'isMaximized':
          return _maximized;
        case 'isMinimized':
          return _minimized;
        case 'setFullScreen':
          _fullscreen = call.arguments['isFullScreen'] as bool;
          return null;
        case 'maximize':
          _maximized = true;
          return null;
        case 'unmaximize':
          _maximized = false;
          return null;
        case 'restore':
          _minimized = false;
          return null;
        default:
          return null;
      }
    },
  );
}

void _clearMock() {
  TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger.setMockMethodCallHandler(
    const MethodChannel('window_manager'),
    null,
  );
}

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  group('WindowManagerX real defaults', () {
    setUp(_mockWindowManager);
    tearDown(_clearMock);

    test('ensureInitialized invokes the plugin', () async {
      await WindowManagerX().ensureInitialized();
      expect(_calls, contains('ensureInitialized'));
    });

    test('setTitleBarStyle invokes the plugin', () async {
      await WindowManagerX().setTitleBarStyle(TitleBarStyle.hidden);
      expect(_calls, contains('setTitleBarStyle'));
    });

    test('setFullScreen invokes the plugin', () async {
      await WindowManagerX().setFullScreen(true);
      expect(_calls, contains('setFullScreen'));
    });

    test('isFullScreen reads through to the plugin', () async {
      _fullscreen = true;
      expect(await WindowManagerX().isFullScreen(), isTrue);
      expect(_calls, contains('isFullScreen'));
    });

    test('maximize/unmaximize/isMaximized invoke the plugin', () async {
      final wm = WindowManagerX();
      await wm.maximize();
      expect(await wm.isMaximized(), isTrue);
      await wm.unmaximize();
      expect(await wm.isMaximized(), isFalse);
      expect(_calls, containsAllInOrder(['maximize', 'isMaximized', 'unmaximize', 'isMaximized']));
    });

    test('close invokes the plugin', () async {
      await WindowManagerX().close();
      expect(_calls, contains('close'));
    });

    // Regression: window_manager's own isFullScreen()/isMaximized() getters
    // are unguarded channel reads (no client-side default) - on Linux
    // isFullScreen in particular can report true before the window is
    // realized. waitUntilReadyToShow is the plugin's sanctioned fix: it
    // reads that (possibly stale) state and force-corrects it. Verify our
    // wrapper actually drives that correction rather than silently trusting
    // whatever the channel returned.
    test('waitUntilReadyToShow corrects a stale fullscreen/maximized report', () async {
      _fullscreen = true;
      _maximized = true;

      await WindowManagerX().waitUntilReadyToShow();

      expect(_calls, contains('waitUntilReadyToShow'));
      expect(_calls, contains('setFullScreen'));
      expect(_calls, contains('unmaximize'));
      expect(_fullscreen, isFalse);
      expect(_maximized, isFalse);
    });

    test('waitUntilReadyToShow leaves an already-settled state alone', () async {
      await WindowManagerX().waitUntilReadyToShow();

      expect(_calls, contains('waitUntilReadyToShow'));
      expect(_calls, isNot(contains('setFullScreen')));
      expect(_calls, isNot(contains('unmaximize')));
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

      await wm.ensureInitialized();
      await wm.waitUntilReadyToShow();
      await wm.setTitleBarStyle(TitleBarStyle.hidden);
      await wm.setFullScreen(true);
      await wm.isFullScreen();
      await wm.maximize();
      await wm.unmaximize();
      await wm.isMaximized();
      await wm.close();

      expect(calls, [
        'ensureInitialized',
        'waitUntilReadyToShow',
        'setTitleBarStyle',
        'setFullScreen',
        'isFullScreen',
        'maximize',
        'unmaximize',
        'isMaximized',
        'close',
      ]);
    });

    test('waitUntilReadyToShow invokes the given callback, like the real plugin', () async {
      var called = false;
      await WindowManagerX.fake().waitUntilReadyToShow(null, () => called = true);
      expect(called, isTrue);
    });
  });
}
