import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import 'package:retrovibed/uuidx.dart' as uuidx;
import './api.dart' as api;
import 'package:retrovibed/design.kit/stateful.dart';

class Settings extends StatefulWidget {
  static api.Wireguard zero = api.Wireguard();
  final api.Wireguard current;
  final Future<api.Wireguard> Function(api.Wireguard)? onChange;

  const Settings(
    this.current, {
    super.key,
    this.onChange,
  });

  static FutureBuilder<api.Wireguard> future(
    Future<api.Wireguard> pending, {
    Future<api.Wireguard> Function(api.Wireguard)? onChange,
  }) {
    return ds.future(Settings.zero, pending, (snapshot) {
      return ds.ErrorScreen(
        Settings(
          snapshot.data ?? Settings.zero,
          key: ValueKey(snapshot.data?.id ?? uuidx.min()),
          onChange: onChange,
        ),
        cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
      );
    });
  }

  @override
  State<Settings> createState() => _SettingsState(this.current);
}

class _SettingsState extends State<Settings> with LoadingState {
  api.Wireguard current;

  _SettingsState(this.current);

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return forms.Container(
      padding: defaults.padding,
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: [
          forms.Field(
            label: Text("port"),
            input: TextFormField(
              decoration: InputDecoration(
                helperText: "0 means automatic port detection",
              ),
              initialValue: current.port.toString(),
              keyboardType: TextInputType.number,
              onChanged: (v) {
                final port = int.tryParse(v);
                if (port == null) return;
                final updated = current..port = port;
                setState(() => current = updated);
                widget.onChange?.call(updated);
              },
            ),
          ),
          forms.Field(
            label: Text("dns rate limit"),
            input: inputs.RateLimit(
              value: current.dnsRateLimit,
              presets: [
                (label: '10/sec', value: 10, unit: 'sec'),
                (label: '100/sec', value: 100, unit: 'sec'),
                (label: '1000/sec', value: 1000, unit: 'sec'),
              ],
              onChanged: (v) {
                final updated = current..dnsRateLimit = v;
                setState(() => current = updated);
                widget.onChange?.call(updated);
              },
            ),
          ),
        ],
      ),
    );
  }
}
