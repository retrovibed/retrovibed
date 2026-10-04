import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/wireguard.dart' as wireguard;

/// Lets the user upload a single wireguard configuration during first time
/// setup. Calls [onDone] once it's uploaded, or when the user continues without one.
class WireguardSetup extends StatefulWidget {
  final VoidCallback onDone;
  const WireguardSetup({super.key, required this.onDone});

  @override
  State<WireguardSetup> createState() => _WireguardSetupState();
}

class _WireguardSetupState extends State<WireguardSetup> with ds.LoadingState {
  @override
  void initState() {
    super.initState();
    loading = false;
  }

  // only the first file is uploaded; additional configs can be added under settings.
  Future<Widget?> _upload(ds.FilesEvent v, {StreamSink<httpx.UploadProgress>? progress}) {
    final auth = [authn.request(authn.AuthzCache.meta(context))];
    final c = v.files.first;
    return wireguard.wireguard
        .uploadable(c.path, c.name, c.mimeType!)
        .then((f) => wireguard.wireguard.upload((req) => req..files.add(f)))
        .then((r) => wireguard.wireguard.touch(r.wireguard.id, r.wireguard.nettype, options: auth))
        .then((_) => widget.onDone())
        .then((_) => ds.NullWidget)
        .catchError((c) {
      setState(() => cause = ds.Error.unknown(c, onTap: reseterr));
      return ds.NullWidget;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      padding: defaults.padding,
      margin: defaults.margin,
      constraints: const BoxConstraints(maxWidth: 512),
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing,
        children: [
          Text("VPN - WireGuard", style: theme.textTheme.titleMedium),
          const Text("Optionally secure your privacy by routing traffic through a WireGuard VPN"),
          Text(
            "You can skip this and set it up later under Settings → VPN - WireGuard.",
            style: theme.textTheme.bodySmall,
          ),
          const Divider(),
          ds.Loading(
            cause: cause,
            loading: loading,
            // fixed height: the drop well needs bounded constraints and the
            // ftux overlay is wrapped in a scroll view.
            SizedBox(
              height: 200,
              child: ds.FileDropWell(
                _upload,
                mimetypes: const [mimex.text.plain],
                child: ds.FileDropWell.textual("drop a wireguard configuration file"),
                shape: RoundedRectangleBorder(borderRadius: defaults.borderRadius),
              ),
            ),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            spacing: defaults.spacing,
            children: [
              ds.LoadingButton(const Text('Continue'), onPressed: () async => widget.onDone()),
            ],
          ),
        ],
      ),
    );
  }
}
