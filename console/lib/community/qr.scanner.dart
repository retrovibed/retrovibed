import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'qr.dart';
import 'qr.camera.dart';
import 'community.pb.dart';
import 'community.detail.dart';

class QRScannerModal extends StatefulWidget {
  final Future<void> Function(Community community, String attribution) onScanned;
  final VoidCallback onCancel;
  final Widget Function(void Function(String data) onDetect, void Function(Object, StackTrace) onError) camera;

  const QRScannerModal({
    super.key,
    required this.onScanned,
    required this.onCancel,
    this.camera = QRCameraView.defaults,
  });

  @override
  State<QRScannerModal> createState() => _QRScannerModalState();
}

class _QRScannerModalState extends State<QRScannerModal> {
  Widget _cause = ds.Error.zero;
  bool _processing = false;

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
      _processing = false;
    });
  }

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _processQRData(String data) {
    if (_processing) return;
    setState(() {
      _processing = true;
      _cause = ds.Error.zero;
    });

    final (community, attribution) = decodeQRPayload(data);
    if (community == null) {
      return setState(() {
        _processing = false;
        _cause = ds.Error.unknown('invalid QR', onTap: _reseterr);
      });
    }

    setState(
      () =>
          _cause = ds.Masked(
            alignment: Alignment.center,
            ds.Confirmation.yesNo(
              content: Column(
                children: [
                  CommunityDetail(
                    community: community,
                  ),
                ],
              ),
              onConfirm: (_) {
                widget.onScanned(community, attribution).catchError((e, s) {
                  setState(() => _cause = ds.Error.unknown(e, onTap: _reseterr));
                  return Future.error(e);
                }).ignore();
              },
              onCancel: (_) => _reseterr(),
            ),
          ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      alignment: Alignment.center,
      constraints: BoxConstraints(maxWidth: 512, maxHeight: 600),
      padding: defaults.padding,
      margin: defaults.margin,
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Flexible(
                child: Text(
                  'Scan QR Code',
                  style: theme.textTheme.headlineSmall,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              IconButton(onPressed: widget.onCancel, icon: Icon(Icons.close)),
            ],
          ),
          Expanded(
            child: ds.ErrorScreen(
              cause: _cause,
              ClipRRect(
                borderRadius: defaults.borderRadius,
                child: widget.camera(
                  _processQRData,
                  (Object e, StackTrace s) {
                    setState(() {
                      _processing = false;
                      _cause = ds.Error.unknown(e, onTap: _reseterr);
                    });
                  },
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

void showQRScanner(
  BuildContext context, {
  required Future<void> Function(Community community, String attribution) onScanned,
}) {
  final modal = ds.modals.of(context);
  ds.modals.push(
    context,
    QRScannerModal(
      onScanned: (community, attribution) => onScanned(community, attribution).then((_) => modal?.push(null)),
      onCancel: () => modal?.push(null),
    ),
  );
}
