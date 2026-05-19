import 'dart:io';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/billing.dart' as billing;
import 'package:retrovibed/quotas.dart' as quotas;
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/rss.dart' as rss;
import 'package:retrovibed/wireguard.dart' as wireguard;
import 'package:retrovibed/google.dart' as google;
import 'package:retrovibed/usermanagement.dart' as usermanagement;
import 'package:retrovibed/debug.dart' as debug;

class AutoHelp extends StatelessWidget {
  final Widget child;
  const AutoHelp(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ds.HelpAuto(
      child,
      cacheid: 'settings',
      title: Text("Settings", style: theme.textTheme.titleMedium),
      content: const Text(
        "Configure your account, billing, storage, and integrations. "
        "Manage connected services, RSS feeds, VPN settings, and "
        "user access from this panel.",
      ),
    );
  }
}

class Display extends StatefulWidget {
  const Display({super.key});

  @override
  State<Display> createState() => _DisplayState();
}

class _DisplayState extends State<Display> {
  Widget _overlay = ds.Empty;
  ValueNotifier<meta.Daemon> _library = ValueNotifier(meta.Daemon());
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void refresh() {
    setState(() {}); // force rebuild
  }

  @override
  void initState() {
    super.initState();
    _library = meta.EndpointAuto.of(context)?.changed ?? _library;
    _library.addListener(refresh);
  }

  @override
  void dispose() {
    super.dispose();
    _library.removeListener(refresh);
  }

  void overlay(Widget w) {
    setState(() {
      _overlay = w;
    });
  }

  void masked(Widget w) {
    overlay(
      ds.Masked(alignment: Alignment.center, reset: () => overlay(ds.Empty), w),
    );
  }

  void full(Widget w) {
    overlay(
      ds.build((context) {
        final defaults = ds.Defaults.of(context);
        return ds.Container(
          alignment: defaults.isCompact ? Alignment.bottomCenter : Alignment.topCenter,
          margin: EdgeInsets.zero,
          padding: EdgeInsets.zero,
          SingleChildScrollView(child: w),
        );
      }),
    );
  }

  @override
  Widget build(BuildContext context) {
    return ds.build((context) {
      final defaults = ds.Defaults.of(context);
      final compact = defaults.isCompact;
      final _billing = billing.Registered.of(context);
      final _displaybilling = !(_billing.current.subscriptionId.isEmpty && Platform.isMacOS);

      return SelectionArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          mainAxisAlignment: MainAxisAlignment.start,
          verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
          children: [
            meta.DaemonDropdown(
              library: _library,
              trailing: [
                _overlay == ds.Empty
                    ? IconButton(
                      onPressed: () {
                        masked(
                          ds.Confirmation.yesNo(
                            content: Text(
                              'Delete ${_library.value.description}?',
                            ),
                            onCancel: () => overlay(ds.Empty),
                            onConfirm: () {
                              httpx.withRetry(
                                () => meta.daemons.delete(_library.value.id).then((_) {
                                  overlay(ds.Empty);
                                }),
                              );
                            },
                          ),
                        );
                      },
                      icon: Icon(Icons.delete),
                    )
                    : ds.LoadingIconButton.close(
                      onPressed: () {
                        overlay(ds.Empty);
                        return Future.value(null);
                      },
                    ),
              ],
            ),
            Expanded(
              child: Padding(
                padding: EdgeInsets.symmetric(
                  horizontal: defaults.padding.horizontal / 2,
                ),
                child: ds.Overlay(
                  alignment: Alignment.topLeft,
                  ds.layout((context, constraints) {
                    const mainAxisExtent = 192.0;
                    const crossAxisExtent = 192.0;
                    int crossAxisCount;

                    if (constraints.maxWidth > mainAxisExtent) {
                      crossAxisCount = constraints.maxWidth ~/ mainAxisExtent;
                    } else {
                      crossAxisCount = 1;
                    }

                    return GridView(
                      reverse: compact,
                      gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                        crossAxisCount: crossAxisCount,
                        mainAxisExtent: mainAxisExtent,
                        childAspectRatio: crossAxisExtent / mainAxisExtent,
                        crossAxisSpacing: defaults.spacing / 2,
                        mainAxisSpacing: defaults.spacing / 2,
                      ),
                      children: [
                        billing.Card(
                          onPressed: full,
                          margin: EdgeInsets.zero,
                        ),
                        if (_displaybilling) ...[
                          billing.ReferralCard(
                            onPressed: full,
                            margin: EdgeInsets.zero,
                          ),
                          billing.InviteCard(margin: EdgeInsets.zero),
                          quotas.Card(),
                        ],
                        profiles.Card(
                          onPressed: defaults.debug ? full : null,
                        ),
                        rss.Card(
                          onPressed: full,
                          margin: EdgeInsets.zero,
                        ),
                        wireguard.Card(
                          onPressed: full,
                          margin: EdgeInsets.zero,
                        ),
                        usermanagement.Card(
                          onPressed: full,
                          margin: EdgeInsets.zero,
                        ),
                        google.Card(onPressed: full),
                        debug.Card(margin: EdgeInsets.zero),
                        if (defaults.debug) debug.MeteredCard(margin: EdgeInsets.zero),
                      ],
                    );
                  }),
                  overlay: MediaQuery(
                    data: MediaQuery.of(context).copyWith(
                      padding: EdgeInsets.only(
                        top: 139,
                      ), // compensate for the dropdown and titlebar
                    ),
                    child: _overlay,
                  ),
                ),
              ),
            ),
          ],
        ),
      );
    });
  }
}
