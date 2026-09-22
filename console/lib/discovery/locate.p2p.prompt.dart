import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'api.dart' as api;
import 'settings.locate.dart';

// P2PConsentDeclined is the rejection error ensureP2P's Future completes
// with when the user declines the disclaimer prompt - callers can match on
// it via consentDeclined to abort the locate/download quietly, distinct
// from a real failure.
class P2PConsentDeclined implements Exception {
  const P2PConsentDeclined();

  @override
  String toString() => "P2PConsentDeclined: user declined p2p discovery consent";
}

bool consentDeclined(Object obj) => obj is P2PConsentDeclined;

// ensureLocateP2P checks whether P2P discovery is enabled server-side; if it
// isn't, it prompts the user with the same legal disclaimer shown in
// Settings, and on confirmation persists locateP2p = true before returning.
// Rejects with P2PConsentDeclined if the user declines.
Future<void> ensureP2P(
  BuildContext context, {
  List<httpx.Option> options = const [],
}) {
  return api.configuration.get(options: options).then((settings) {
    if (settings.locateP2p) return Future<void>.value();

    return ds.modals
        .asyncfn<void>(
          context,
          (completion) => ds.Confirmation.yesNo(
            content: const Text(LocateSettings.disclaimerText),
            onConfirm: (_) => completion.complete(),
            onCancel: (_) => completion.completeError(const P2PConsentDeclined()),
          ),
        )
        .then<void>((_) => api.configuration.create(settings..locateP2p = true, options: options));
  });
}
