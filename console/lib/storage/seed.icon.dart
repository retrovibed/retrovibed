import 'package:flutter/material.dart';
import './seed.dart' as seed;

class SeedIcon extends StatelessWidget {
  final String current;
  final seed.Classifier classifier;
  final double size;

  const SeedIcon(
    this.current, {
    super.key,
    this.size = 24.0,
    required this.classifier,
  });

  @override
  Widget build(BuildContext context) {
    final classifer = this.classifier;
    final _current = classifer.classify(this.current);
    return Tooltip(
      message: _current.tooltip,
      child: Icon(_current.icon, size: size),
    );
  }
}
