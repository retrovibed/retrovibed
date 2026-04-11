import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/community/link.content.dart';
import 'package:retrovibed/community/qr.scanner.dart';
import 'package:retrovibed/community/list.display.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

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
                      help: ds.Hint(
                        label: const Text("QR"),
                        description: const Text(
                          "scan a QR code to subscribe or link content",
                        ),
                      ),
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
