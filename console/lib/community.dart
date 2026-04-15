import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/community/link.content.dart';
import 'package:retrovibed/community/qr.scanner.dart';
import 'package:retrovibed/community/list.display.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

class AutoHelp extends StatelessWidget {
  final Widget child;
  const AutoHelp(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ds.HelpAuto(
      child,
      cacheid: 'community',
      title: Text("Community", style: theme.textTheme.titleMedium),
      content: const Text(
        "Connect and share with other users. Browse community-curated "
        "collections, join groups, and discover new content shared by "
        "people with similar tastes.\n\n"
        "Press Alt+? at any time to activate/deactivate help overlay",
      ),
    );
  }
}

class Management extends StatefulWidget {
  @override
  _ManagementState createState() => _ManagementState();
}

class _ManagementState extends State<Management> {
  ValueKey<String> _refresh = ValueKey(uuidx.random());
  Widget _overlay = ds.Empty;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _resetOverlay() {
    setState(() {
      _overlay = ds.Empty;
    });
  }

  @override
  Widget build(BuildContext context) {
    return ds.build((context) {
      final theme = Theme.of(context);
      final defaults = ds.Defaults.of(context);
      final compact = defaults.isCompact;
      return ds.Overlay(
        ds.Container(
          margin: EdgeInsets.zero,
          padding: defaults.padding,
          Column(
            mainAxisSize: MainAxisSize.max,
            mainAxisAlignment: MainAxisAlignment.start,
            verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                spacing: defaults.spacing,
                children: [
                  Expanded(
                    child: Text(
                      'Communities',
                      style: theme.textTheme.headlineMedium,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  if (!defaults.desktop)
                    ds.LoadingIconButton(
                      icon: Icon(Icons.qr_code_scanner),
                      onPressed: () {
                        setState(() {
                          _overlay = QRScannerModal(
                            onScanned:
                                (community, attribution) => handleSubscribeAction(context, community, attribution).then(
                                  (_) => setState(() {
                                    _overlay = ds.Empty;
                                  }),
                                ),
                            onCancel: _resetOverlay,
                          );
                        });
                        return Future.value();
                      },
                      help: ds.Hint(const Text("scan a QR code to subscribe or link content")),
                    ),
                ],
              ),
              Expanded(child: ListDisplay(key: _refresh)),
            ],
          ),
        ),
        overlay: _overlay,
      );
    });
  }
}
