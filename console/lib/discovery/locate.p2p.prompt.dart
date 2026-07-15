import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'api.dart' as api;
import 'settings.locate.dart';

// ensureLocateP2P checks whether P2P discovery is enabled server-side; if it
// isn't, it prompts the user with the same legal disclaimer shown in
// Settings, and on confirmation persists locateP2p = true before returning.
// Returns false (no error) if the user declines - callers should abort the
// locate quietly in that case.
Future<bool> ensureP2P(
  BuildContext context, {
  List<httpx.Option> options = const [],
}) async {
  final settings = await api.configuration.get(options: options);
  if (settings.locateP2p) return true;

  final proceed = await ds.modals.asyncfn<bool>(
    context,
    (completion) => ds.Confirmation.yesNo(
      content: const Text(LocateSettings.disclaimerText),
      onConfirm: (_) => completion.complete(true),
      onCancel: (_) => completion.complete(false),
    ),
  );
  if (proceed != true) return false;

  await api.configuration.create(settings..locateP2p = true, options: options);
  return true;
}
