import 'package:material_ui/material_ui.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;

class Settings extends StatefulWidget {
  static api.DHTSettings zero = api.DHTSettings(
    passive: false,
    autoBootstrap: false,
    log: false,
    debug: false,
  );
  final api.DHTSettings defaults;
  final Future<api.DHTSettings> Function(api.DHTSettings)? onChange;

  Settings(this.defaults, {super.key, this.onChange});

  static FutureBuilder<api.DHTSettings> future(
    Future<api.DHTSettings> pending, {
    Future<api.DHTSettings> Function(api.DHTSettings)? onChange,
  }) {
    return ds.future(Settings.zero, pending, (snapshot) {
      return ds.ErrorScreen(
        Settings(
          snapshot.data ?? Settings.zero,
          key: ValueKey(snapshot.data.hashCode),
          onChange: onChange,
        ),
        cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
      );
    });
  }

  @override
  State<Settings> createState() => _EditView(this.defaults);
}

class _EditView extends State<Settings> {
  api.DHTSettings current;
  final ValueNotifier<api.DHTSettings> _update = ValueNotifier<api.DHTSettings>(api.DHTSettings());

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
          Row(
            spacing: defaults.spacing,
            children: [
              Expanded(
                child: Tooltip(
                  message:
                      "only answer dht queries instead of actively participating in the distributed hash table to find peers",
                  child: forms.Checkbox(
                    const Text("passive"),
                    value: current.passive,
                    onChanged: (v) {
                      final _update = current.deepCopy()..passive = v ?? current.passive;
                      setState(() => current = _update);
                    },
                  ),
                ),
              ),
              Expanded(
                child: Tooltip(
                  message: "bootstrap the dht from the public global bootstrap nodes",
                  child: forms.Checkbox(
                    const Text("bootstrap"),
                    value: current.autoBootstrap,
                    onChanged: (v) {
                      final _update = current.deepCopy()..autoBootstrap = v ?? current.autoBootstrap;
                      setState(() => current = _update);
                    },
                  ),
                ),
              ),
            ],
          ),
          Row(
            spacing: defaults.spacing,
            children: [
              Text("logging", style: Theme.of(context).textTheme.labelMedium),
              const Expanded(child: Divider()),
            ],
          ),
          Row(
            spacing: defaults.spacing,
            children: [
              Expanded(
                child: Tooltip(
                  message: "log dht activity",
                  child: forms.Checkbox(
                    const Text("info"),
                    value: current.log,
                    onChanged: (v) {
                      final _update = current.deepCopy()..log = v ?? current.log;
                      setState(() => current = _update);
                    },
                  ),
                ),
              ),
              Expanded(
                child: Tooltip(
                  message: "log every dht message",
                  child: forms.Checkbox(
                    const Text("debug"),
                    value: current.debug,
                    onChanged: (v) {
                      final _update = current.deepCopy()..debug = v ?? current.debug;
                      setState(() => current = _update);
                    },
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
