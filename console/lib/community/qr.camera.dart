import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:retrovibed/designkit.dart' as ds;

class QRCameraView extends StatefulWidget {
  final void Function(String data) onDetect;
  final void Function(Object, StackTrace) onError;
  const QRCameraView({
    super.key,
    required this.onDetect,
    required this.onError,
  });

  static Widget defaults(void Function(String) onDetect, void Function(Object, StackTrace) onError) =>
      QRCameraView(onDetect: onDetect, onError: onError);

  @override
  State<QRCameraView> createState() => _QRCameraViewState();
}

class _QRCameraViewState extends State<QRCameraView> {
  final MobileScannerController _controller = MobileScannerController();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MobileScanner(
      controller: _controller,
      onDetect: (capture) {
        final value = capture.barcodes.firstOrNull?.rawValue;
        if (value != null) widget.onDetect(value);
      },
      onDetectError: widget.onError,
      errorBuilder: (context, error) {
        return ds.Error.unknown(error, onTap: () => widget.onError(error, StackTrace.current));
      },
    );
  }
}
