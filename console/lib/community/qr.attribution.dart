import 'package:flutter/material.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/qr.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class QRAttribution extends StatefulWidget {
  final Community community;

  const QRAttribution({super.key, required this.community});

  @override
  State<QRAttribution> createState() => _QRAttributionState();
}

class _QRAttributionState extends State<QRAttribution> with LoadingState {
  String _attribution = '';

  @override
  void initState() {
    super.initState();
    authn.DeeppoolAuthzCache.attributionToken(context).then((v) {
      setState(() => _attribution = v);
    }).ignore();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final qrData = encodeQRPayload(widget.community, attribution: _attribution);

    return ds.Help(
      ds.Container(
        clipBehavior: Clip.antiAlias,
        constraints: BoxConstraints(maxHeight: defaults.compact + defaults.padding.vertical),
        ClipRRect(
          borderRadius: defaults.borderRadius,
          child: QrImageView(
            data: qrData,
            version: QrVersions.auto,
            backgroundColor: Colors.white,
            dataModuleStyle: QrDataModuleStyle(color: Colors.black),
          ),
        ),
      ),
      ds.Hint(Text("Scanning this QR code will subscribe users to this community")),
    );
  }
}
