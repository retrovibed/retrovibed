import 'package:flutter/material.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/qr.dart';

class QRAttribution extends StatefulWidget {
  final Community community;

  const QRAttribution({super.key, required this.community});

  @override
  State<QRAttribution> createState() => _QRAttributionState();
}

class _QRAttributionState extends State<QRAttribution> {
  String _attribution = '';

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

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

    return ds.Container(
      margin: defaults.margin.copyWith(left: 0.0, right: 0.0),
      clipBehavior: Clip.antiAlias,
      constraints: BoxConstraints(maxHeight: defaults.compact, maxWidth: defaults.compact),
      ClipRRect(
        borderRadius: defaults.borderRadius,
        child: QrImageView(
          data: qrData,
          version: QrVersions.auto,
          backgroundColor: Colors.white,
          dataModuleStyle: QrDataModuleStyle(color: Colors.black),
        ),
      ),
    );
  }
}
