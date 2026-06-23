import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;

const String _disclaimerCacheId = 'discovery.p2p';

const String _disclaimerText =
    'Retrovibed operates on an entirely peer-to-peer (p2p) system.\n\n'
    'The p2p system allows journalists, musicians, and other content creators to directly '
    'interact with users and be discovered, preventing censorship by third parties '
    '(including retrovibed).\n\n'
    'As a result, it can locate media that has been published by third parties who do not '
    'have distribution rights. Retrovibed takes no responsibility for such content.\n\n'
    'P2P: operate in fully p2p mode, allowing discovery of publicly available content. '
    'You take responsibility for your activities and for obeying all legal laws of your region.';

class LocateSettings extends StatefulWidget {
  static api.DiscoverySettings zero = api.DiscoverySettings(locateP2p: false);
  final api.DiscoverySettings defaults;
  final Future<api.DiscoverySettings> Function(api.DiscoverySettings)? onChange;
  final bool Function(String)? disclaimer;
  final void Function(String)? acknowledge;

  LocateSettings(
    this.defaults, {
    super.key,
    this.onChange,
    this.disclaimer,
    this.acknowledge,
  });

  static FutureBuilder<api.DiscoverySettings> future(
    Future<api.DiscoverySettings> pending, {
    Future<api.DiscoverySettings> Function(api.DiscoverySettings)? onChange,
  }) {
    return ds.future(LocateSettings.zero, pending, (snapshot) {
      return ds.ErrorScreen(
        LocateSettings(
          snapshot.data ?? LocateSettings.zero,
          key: ValueKey(snapshot.data.hashCode),
          onChange: onChange,
        ),
        cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
      );
    });
  }

  @override
  State<LocateSettings> createState() => _LocateEditView(this.defaults);
}

class _LocateEditView extends State<LocateSettings> {
  api.DiscoverySettings current;

  _LocateEditView(this.current);

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _update(api.DiscoverySettings updated) {
    setState(() => current = updated);
    widget.onChange?.call(current);
  }

  @override
  Widget build(BuildContext context) {
    return ds.DisclaimerIntercept(
      forms.Checkbox(
        const Text("p2p"),
        dense: true,
        value: current.locateP2p,
        help: ds.Hint(
          const Text("locate media via the distributed p2p discovery network, disabled by default"),
        ),
        onChanged: (v) => _update(current..locateP2p = v ?? !current.locateP2p),
      ),
      cacheid: _disclaimerCacheId,
      cached: widget.disclaimer ?? ds.Disclaimer.disclaimerpath,
      acknowledge: widget.acknowledge ?? ds.Disclaimer.acknowledge,
      overlay: (complete) => ds.Confirmation.yesNo(
        content: const Text(_disclaimerText),
        onConfirm: (_) {
          complete(true);
          _update(current..locateP2p = !current.locateP2p);
        },
        onCancel: (_) => complete(false),
      ),
    );
  }
}
