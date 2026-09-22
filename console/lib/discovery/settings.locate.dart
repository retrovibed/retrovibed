import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/stateful.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;

const String _disclaimerCacheId = 'discovery.p2p';

class LocateSettings extends StatefulWidget {
  static api.DiscoverySettings zero = api.DiscoverySettings(locateP2p: false);

  static const String disclaimerText = '''
Retrovibed supports an entirely peer-to-peer (p2p) environment.

The p2p system allows journalists, musicians, and other content creators to
directly interact with users and be discovered, preventing censorship by
third parties (including retrovibed).

As a result, it can locate media that has been published by parties who do not
have distribution rights. Such content is not a part of Retrovibed's platform and
we takes no responsibility for such content.

By enabling P2P discovery of available content. You take responsibility for your
activities and for obeying the laws within your region.''';
  final api.DiscoverySettings defaults;
  final Future<api.DiscoverySettings> Function(api.DiscoverySettings) onChange;
  final bool Function(String)? disclaimer;
  final void Function(String)? acknowledge;

  LocateSettings(
    this.defaults, {
    super.key,
    this.onChange = ds.fnAsyncPassthrough,
    this.disclaimer,
    this.acknowledge,
  });

  static FutureBuilder<api.DiscoverySettings> future(
    Future<api.DiscoverySettings> pending, {
    Future<api.DiscoverySettings> Function(api.DiscoverySettings) onChange = ds.fnAsyncPassthrough,
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
  State<LocateSettings> createState() => _LocateEditView();
}

class _LocateEditView extends State<LocateSettings> with LoadingState {
  api.DiscoverySettings current = api.DiscoverySettings();

  @override
  void initState() {
    super.initState();
    loading = false;
    current = widget.defaults;
  }

  void _update(api.DiscoverySettings updated) {
    setState(() => loading = true);

    Future.sync(() => widget.onChange(updated))
        .then((v) {
          setState(() {
            current = v;
            loading = false;
          });
        })
        .catchError((error) {
          setState(() {
            loading = false;
            cause = ds.Errors.httpauto(error, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((error) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(error, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return forms.Container(
      cause: cause,
      loading: loading,
      ds.DisclaimerIntercept(
        forms.Checkbox(
          const Text("p2p"),
          dense: true,
          value: current.locateP2p,
          help: ds.Hint(
            const Text("locate media via the distributed p2p discovery network, disabled by default"),
          ),
          onChanged: loading ? null : (v) => _update(current..locateP2p = v ?? !current.locateP2p),
        ),
        cacheid: _disclaimerCacheId,
        cached: widget.disclaimer ?? ds.Disclaimer.disclaimerpath,
        acknowledge: widget.acknowledge ?? ds.Disclaimer.acknowledge,
        overlay: (complete) => ds.Confirmation.yesNo(
          content: const Text(LocateSettings.disclaimerText),
          onConfirm: (_) {
            complete(true);
            _update(current..locateP2p = !current.locateP2p);
          },
          onCancel: (_) => complete(false),
        ),
      ),
    );
  }
}
