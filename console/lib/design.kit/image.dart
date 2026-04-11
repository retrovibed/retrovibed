import 'package:flutter/material.dart' as m;

class Image extends m.StatelessWidget {
  final String current;
  final double? size;
  const Image(
    this.current, {
    super.key,
    this.size,
  });

  @override
  m.Widget build(m.BuildContext context) {
    return current == "" ? m.Icon(m.Icons.image_outlined, size: size) : m.Image.network(current, height: size);
  }
}
