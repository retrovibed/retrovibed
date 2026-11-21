import 'package:material_ui/material_ui.dart';
import 'package:retrovibed/meta/api.dart' as api;

// overridable seam over the static api.daemons, whose implementations reach
// into the native library via FFI. widget tests swap this out through
// DaemonsScope so they never need to load libretrovibed.so.
class RetroAPI {
  final bool Function(api.Daemon) isLocalDevice;

  const RetroAPI({
    this.isLocalDevice = api.daemons.isLocalDevice,
  });
}

class RetroFFIScope extends InheritedWidget {
  final RetroAPI daemons;

  const RetroFFIScope({
    super.key,
    this.daemons = const RetroAPI(),
    required super.child,
  });

  // falls back to the real implementation when no scope is mounted, so
  // production needs no wiring.
  static RetroAPI of(BuildContext context) =>
      context.dependOnInheritedWidgetOfExactType<RetroFFIScope>()?.daemons ?? const RetroAPI();

  @override
  bool updateShouldNotify(RetroFFIScope old) => daemons != old.daemons;
}

extension RetroFFIContext on BuildContext {
  RetroAPI get ffiretro => RetroFFIScope.of(this);
}
