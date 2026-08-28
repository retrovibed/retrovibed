import 'package:flutter/foundation.dart';
import 'package:window_manager/window_manager.dart';

typedef FnVoid = Future<void> Function();
typedef FnBool = Future<bool> Function();
typedef FnSetFullScreen = Future<void> Function(bool);
typedef FnSetTitleBarStyle = Future<void> Function(TitleBarStyle, {bool windowButtonVisibility});
typedef FnWaitUntilReadyToShow = Future<void> Function([WindowOptions?, VoidCallback?]);
typedef FnAddListener = void Function(WindowListener);
typedef FnRemoveListener = void Function(WindowListener);

// WindowManagerX bundles the window_manager plugin calls the app actually
// uses behind injectable functions so widgets/bootstrap can be tested
// without a real platform channel. Every function defaults to the real
// windowManager singleton's method.
class WindowManagerX {
  final FnVoid ensureInitialized;
  final FnWaitUntilReadyToShow waitUntilReadyToShow;
  final FnSetTitleBarStyle setTitleBarStyle;
  final FnSetFullScreen setFullScreen;
  final FnBool isFullScreen;
  final FnVoid maximize;
  final FnVoid unmaximize;
  final FnBool isMaximized;
  final FnVoid close;
  final FnAddListener addListener;
  final FnRemoveListener removeListener;

  WindowManagerX({
    FnVoid? ensureInitialized,
    FnWaitUntilReadyToShow? waitUntilReadyToShow,
    FnSetTitleBarStyle? setTitleBarStyle,
    FnSetFullScreen? setFullScreen,
    FnBool? isFullScreen,
    FnVoid? maximize,
    FnVoid? unmaximize,
    FnBool? isMaximized,
    FnVoid? close,
    FnAddListener? addListener,
    FnRemoveListener? removeListener,
  }) : ensureInitialized = ensureInitialized ?? windowManager.ensureInitialized,
       waitUntilReadyToShow = waitUntilReadyToShow ?? windowManager.waitUntilReadyToShow,
       setTitleBarStyle = setTitleBarStyle ?? windowManager.setTitleBarStyle,
       setFullScreen = setFullScreen ?? windowManager.setFullScreen,
       isFullScreen = isFullScreen ?? windowManager.isFullScreen,
       maximize = maximize ?? windowManager.maximize,
       unmaximize = unmaximize ?? windowManager.unmaximize,
       isMaximized = isMaximized ?? windowManager.isMaximized,
       close = close ?? windowManager.close,
       addListener = addListener ?? windowManager.addListener,
       removeListener = removeListener ?? windowManager.removeListener;

  // fake returns a WindowManagerX whose state getters answer with the given
  // fullscreen/maximized values (explicitly false by default - never
  // "whatever the plugin happens to report") and whose calls are recorded,
  // in invocation order, into `calls` when provided. Lets a test simulate
  // the plugin reporting a stale/incorrect state (e.g. fullscreen: true
  // before the window is actually fullscreen) without touching a channel.
  static WindowManagerX fake({bool fullscreen = false, bool maximized = false, List<String>? calls}) {
    void record(String method) => calls?.add(method);

    return WindowManagerX(
      ensureInitialized: () async => record('ensureInitialized'),
      waitUntilReadyToShow: ([options, callback]) async {
        record('waitUntilReadyToShow');
        callback?.call();
      },
      setTitleBarStyle: (style, {windowButtonVisibility = true}) async => record('setTitleBarStyle'),
      setFullScreen: (v) async => record('setFullScreen'),
      isFullScreen: () async {
        record('isFullScreen');
        return fullscreen;
      },
      maximize: () async => record('maximize'),
      unmaximize: () async => record('unmaximize'),
      isMaximized: () async {
        record('isMaximized');
        return maximized;
      },
      close: () async => record('close'),
      addListener: (l) => record('addListener'),
      removeListener: (l) => record('removeListener'),
    );
  }
}

final windowx = WindowManagerX();
