import 'package:material_ui/material_ui.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'ffi.retro.api.dart';
import 'api.dart' as api;

class DaemonTypography extends StatelessWidget {
  final api.Daemon current;
  final List<Widget> leading;
  final List<Widget> trailing;

  const DaemonTypography(
    this.current, {
    super.key,
    this.leading = const [],
    this.trailing = const [],
  });

  static String description(api.Daemon current, {required bool Function(api.Daemon) isLocalDevice}) {
    return isLocalDevice(current) ? "local library" : current.description;
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final defaultDisplay = description(current, isLocalDevice: context.ffiretro.isLocalDevice);

    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      spacing: defaults.spacing,
      children: [
        ...leading,
        Flexible(child: Text(defaultDisplay, overflow: TextOverflow.ellipsis)),
        ...trailing,
      ],
    );
  }
}
