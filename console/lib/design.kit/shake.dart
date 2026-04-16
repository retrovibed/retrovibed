import 'package:flutter/material.dart';
import 'package:shake/shake.dart' as _shake;

class ShakeDetector extends StatefulWidget {
  final Widget child;
  final VoidCallback? onShake;
  final double shakeThresholdGravity;
  final int shakeSlopTimeMS;
  final int shakeCountResetTime;
  final int minimumShakeCount;
  final bool useFilter;

  const ShakeDetector({
    super.key,
    required this.child,
    this.onShake,
    this.shakeThresholdGravity = 2.7,
    this.shakeSlopTimeMS = 500,
    this.shakeCountResetTime = 3000,
    this.minimumShakeCount = 1,
    this.useFilter = false,
  });

  @override
  State<ShakeDetector> createState() => _ShakeDetectorState();
}

class _ShakeDetectorState extends State<ShakeDetector> {
  _shake.ShakeDetector? _detector;

  @override
  void initState() {
    super.initState();
    if (widget.onShake != null) {
      _detector = _shake.ShakeDetector.autoStart(
        onPhoneShake: (_) => widget.onShake!(),
        shakeThresholdGravity: widget.shakeThresholdGravity,
        shakeSlopTimeMS: widget.shakeSlopTimeMS,
        shakeCountResetTime: widget.shakeCountResetTime,
        minimumShakeCount: widget.minimumShakeCount,
        useFilter: widget.useFilter,
      );
    }
  }

  @override
  void dispose() {
    _detector?.stopListening();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
