import 'dart:io' show exit;

import 'package:flutter/foundation.dart';
import 'package:window_manager/window_manager.dart' as wm;

typedef FnVoid = Future<void> Function();
typedef FnBool = Future<bool> Function();
typedef FnSetFullScreen = Future<void> Function(bool);
typedef FnSetTitleBarStyle = Future<void> Function(wm.TitleBarStyle, {bool windowButtonVisibility});
typedef FnAddListener = void Function(WindowEventsX);
typedef FnRemoveListener = void Function(WindowEventsX);

// WindowEventsX is the set of window state changes the app reacts to. Pass
// one to WindowManagerX.addListener and the same instance to removeListener.
class WindowEventsX {
  final VoidCallback? onEnterFullScreen;
  final VoidCallback? onLeaveFullScreen;
  final VoidCallback? onMaximize;
  final VoidCallback? onUnmaximize;

  const WindowEventsX({this.onEnterFullScreen, this.onLeaveFullScreen, this.onMaximize, this.onUnmaximize});
}

// WindowManagerX bundles the window calls the app actually uses behind
// injectable functions so widgets/bootstrap can be tested without a real
// window. Every function defaults to the current native window
// (window_manager's native API, built on nativeapi).
class WindowManagerX {
  // settle force-corrects fullscreen/maximized/minimized to a known
  // (windowed) state before anything else reads them.
  final FnVoid settle;
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
    FnVoid? settle,
    FnSetTitleBarStyle? setTitleBarStyle,
    FnSetFullScreen? setFullScreen,
    FnBool? isFullScreen,
    FnVoid? maximize,
    FnVoid? unmaximize,
    FnBool? isMaximized,
    FnVoid? close,
    FnAddListener? addListener,
    FnRemoveListener? removeListener,
  }) : settle = settle ?? _native.settle,
       setTitleBarStyle = setTitleBarStyle ?? _native.setTitleBarStyle,
       setFullScreen = setFullScreen ?? _native.setFullScreen,
       isFullScreen = isFullScreen ?? _native.isFullScreen,
       maximize = maximize ?? _native.maximize,
       unmaximize = unmaximize ?? _native.unmaximize,
       isMaximized = isMaximized ?? _native.isMaximized,
       close = close ?? _native.close,
       addListener = addListener ?? _native.addListener,
       removeListener = removeListener ?? _native.removeListener;

  // fake returns a WindowManagerX whose state getters answer with the given
  // fullscreen/maximized values (explicitly false by default - never
  // "whatever the platform happens to report") and whose calls are recorded,
  // in invocation order, into `calls` when provided. Lets a test simulate
  // the platform reporting a stale/incorrect state (e.g. fullscreen: true
  // before the window is actually fullscreen) without a native window.
  static WindowManagerX fake({bool fullscreen = false, bool maximized = false, List<String>? calls}) {
    void record(String method) => calls?.add(method);

    return WindowManagerX(
      settle: () async => record('settle'),
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

final _native = _NativeWindow();

// _NativeWindow drives the app's current native window. Every call degrades
// to a no-op (or false) when there is no window or the native library is not
// loaded - as in a test VM - mirroring window_manager's own internal guard.
class _NativeWindow {
  wm.Window? _cached;
  final Map<WindowEventsX, wm.ListenerId> _listeners = {};

  wm.Window? get _window => _cached ??= _guard(() => wm.WindowManager.instance.getCurrent(), null);

  static T _guard<T>(T Function() fn, T fallback) {
    try {
      return fn();
    } on ArgumentError {
      // native symbol lookup failed: the library is not loaded.
      return fallback;
    }
  }

  T _with<T>(T Function(wm.Window w) fn, T fallback) {
    final w = _window;
    if (w == null) return fallback;
    return _guard(() => fn(w), fallback);
  }

  Future<void> settle() async {
    _with((w) {
      if (w.isFullScreen) w.isFullScreen = false;
      if (w.isMaximized) w.unmaximize();
      if (w.isMinimized) w.restore();
    }, null);
  }

  Future<void> setTitleBarStyle(wm.TitleBarStyle style, {bool windowButtonVisibility = true}) async {
    _with((w) {
      w.titleBarStyle = style;
      w.isWindowControlButtonsVisible = windowButtonVisibility;
    }, null);
  }

  Future<void> setFullScreen(bool v) async {
    _with((w) {
      w.isFullScreen = v;
    }, null);
  }

  Future<bool> isFullScreen() async => _with((w) => w.isFullScreen, false);
  Future<void> maximize() async => _with((w) => w.maximize(), null);
  Future<void> unmaximize() async => _with((w) => w.unmaximize(), null);
  Future<bool> isMaximized() async => _with((w) => w.isMaximized, false);
  // where the platform can't close a window natively, quit the app instead,
  // as window_manager's 0.5.x close() did.
  Future<void> close() async {
    _with((w) {
      if (!wm.Window.isCloseSupported()) exit(0);
      w.close();
    }, null);
  }

  void addListener(WindowEventsX l) {
    if (_listeners.containsKey(l)) return;
    final id = _with<wm.ListenerId?>(
      (w) => w.addListener((event) {
        switch (event) {
          case wm.WindowEnteredFullScreenEvent():
            l.onEnterFullScreen?.call();
          case wm.WindowExitedFullScreenEvent():
            l.onLeaveFullScreen?.call();
          case wm.WindowMaximizedEvent():
            l.onMaximize?.call();
          // nativeapi has no unmaximize event: leaving the maximized state
          // is reported as a restore, so check whether we're still maximized.
          case wm.WindowRestoredEvent():
            if (!_with((w) => w.isMaximized, false)) l.onUnmaximize?.call();
          default:
            break;
        }
      }),
      null,
    );
    if (id != null) _listeners[l] = id;
  }

  void removeListener(WindowEventsX l) {
    final id = _listeners.remove(l);
    if (id == null) return;
    _with((w) => w.removeListener(id), false);
  }
}

final windowx = WindowManagerX();
