import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;

class LocalStorageSettings extends StatefulWidget {
  final api.Local initial;
  final Future<api.Local> Function(api.Local)? onChange;
  LocalStorageSettings(this.initial, {super.key, this.onChange});

  @override
  State<LocalStorageSettings> createState() => _LocalStorageSettings(this.initial);
}

class _LocalStorageSettings extends State<LocalStorageSettings> {
  api.Local current;
  final ValueNotifier<api.Local> _update = ValueNotifier<api.Local>(
    api.Local(),
  );
  _LocalStorageSettings(this.current);

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
    _update.value = current;
  }

  void handleChange() {
    widget.onChange?.call(this.current);
  }

  @override
  void initState() {
    super.initState();
    _update.addListener(this.handleChange);
  }

  @override
  void dispose() {
    _update.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return forms.Container(
      padding: defaults.padding,
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          forms.Field(
            input: Tooltip(
              message: "maximum storage usage allowed",
              child: inputs.Bytes(
                value: current.maximum.toInt(),
                magnitude: ds.bytesx.GiB,
                decoration: InputDecoration(helperText: "maximum"),
                onChange: (v) {
                  final _update = current.deepCopy()..maximum = v;
                  setState(() => current = _update);
                },
              ),
            ),
          ),
        ],
      ),
    );
  }
}
