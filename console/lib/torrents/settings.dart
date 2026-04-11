import 'dart:math';
import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;

class Settings extends StatefulWidget {
  static api.TorrentSettings zero = api.TorrentSettings(
    seed: true,
    pex: true,
    log: false,
    debug: false,
  );
  final api.TorrentSettings defaults;
  final Future<api.TorrentSettings> Function(api.TorrentSettings)? onChange;

  Settings(this.defaults, {super.key, this.onChange});

  static FutureBuilder<api.TorrentSettings> future(
    Future<api.TorrentSettings> pending, {
    Future<api.TorrentSettings> Function(api.TorrentSettings)? onChange,
  }) {
    return ds.future(Settings.zero, pending, (snapshot) {
      return ds.ErrorScreen(
        Settings(
          snapshot.data ?? Settings.zero,
          key: ValueKey(snapshot.data.hashCode),
          onChange: onChange,
        ),
        cause:
            snapshot.hasError
                ? ds.Error.unknown(snapshot.error!)
                : ds.Error.zero,
      );
    });
  }

  @override
  State<Settings> createState() => _EditView(this.defaults);
}

class _EditView extends State<Settings> {
  api.TorrentSettings current;
  final ValueNotifier<api.TorrentSettings> _update =
      ValueNotifier<api.TorrentSettings>(api.TorrentSettings());

  _EditView(this.current);

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
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        spacing: defaults.spacing,
        children: [
          forms.Field(
            label: Text("download rate"),
            input: Tooltip(
              message:
                  "maximum download rate allowed per second across all downloads, default (0) is unlimited",
              child: inputs.Bytes(
                decoration: new InputDecoration(hintText: "0"),
                magnitude: ds.bytesx.MiB,
                value: current.download.rate.toInt(),
                onChange: (v) {
                  final _update =
                      current.deepCopy()
                        ..download = api.Limit(
                          rate: v.toInt(),
                          burst: current.download.burst,
                        );
                  setState(() => current = _update);
                },
              ),
            ),
          ),
          forms.Field(
            label: Text("upload rate"),
            input: Tooltip(
              message:
                  "maximum upload rate allowed per second across all content, default (0) is unlimited",
              child: inputs.Bytes(
                decoration: new InputDecoration(hintText: "0"),
                value: current.upload.rate.toInt(),
                magnitude: ds.bytesx.MiB,
                onChange: (v) {
                  final _update =
                      current.deepCopy()
                        ..upload = api.Limit(
                          rate: v.toInt(),
                          burst: current.upload.burst,
                        );
                  setState(() => current = _update);
                },
              ),
            ),
          ),
          forms.Field(
            label: Text("peers"),
            input: TextFormField(
              decoration: new InputDecoration(
                hintText: "8",
                helperText: "minimum number of peers wanted when downloading",
              ),
              initialValue: current.peers.min.toString(),
              keyboardType: TextInputType.number,
              onChanged: (v) {
                final _update =
                    current.deepCopy()
                      ..peers = api.Peers(
                        min: int.tryParse(v) ?? current.peers.min,
                        max: current.peers.max,
                      );
                setState(() => current = _update);
              },
            ),
          ),
          forms.Field(
            input: TextFormField(
              decoration: new InputDecoration(
                hintText: "32",
                helperText: "maximum number of peers wanted when downloading",
              ),
              keyboardType: TextInputType.number,
              initialValue: current.peers.max.toString(),
              onChanged: (v) {
                final _update =
                    current.deepCopy()
                      ..peers = api.Peers(
                        min: current.peers.min,
                        max: int.tryParse(v) ?? current.peers.max,
                      );
                setState(() => current = _update);
              },
            ),
          ),
          forms.Field(
            input: TextFormField(
              decoration: new InputDecoration(
                hintText: "32",
                helperText:
                    "maximum number of outbound connections allowed per second",
              ),
              keyboardType: TextInputType.number,
              initialValue: current.outbound.rate.toString(),
              onChanged: (v) {
                final _update =
                    current.deepCopy()
                      ..outbound = api.Limit(
                        rate: int.tryParse(v) ?? current.outbound.rate,
                        burst: max(current.outbound.burst, 1),
                      );
                setState(() => current = _update);
              },
            ),
          ),
          forms.Field(
            input: TextFormField(
              decoration: new InputDecoration(
                hintText: "32",
                helperText:
                    "maximum number of inbound connections allowed per second",
              ),
              keyboardType: TextInputType.number,
              initialValue: current.inbound.rate.toString(),
              onChanged: (v) {
                final _update =
                    current.deepCopy()
                      ..inbound = api.Limit(
                        rate: int.tryParse(v) ?? current.inbound.rate,
                        burst: max(current.inbound.burst, 1),
                      );
                setState(() => current = _update);
              },
            ),
          ),
          ds.Container(
            Wrap(
              children: [
                forms.Checkbox(
                  const Text("log"),
                  value: current.log,
                  onChanged: (v) {
                    final _update = current.deepCopy()..log = v ?? current.log;
                    setState(() => current = _update);
                  },
                ),
                forms.Checkbox(
                  const Text("debug"),
                  value: current.debug,
                  onChanged: (v) {
                    final _update =
                        current.deepCopy()..debug = v ?? current.debug;
                    setState(() => current = _update);
                  },
                ),
                forms.Checkbox(
                  const Text("firewall"),
                  value: current.firewalled,
                  onChanged: (v) {
                    final _update =
                        current.deepCopy()
                          ..firewalled = v ?? current.firewalled;
                    setState(() => current = _update);
                  },
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
