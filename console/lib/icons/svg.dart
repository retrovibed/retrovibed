import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

class SvgIcon extends StatelessWidget {
  final String svg;
  final double? size;
  final Color? color;

  const SvgIcon(this.svg, {super.key, this.size, this.color});

  @override
  Widget build(BuildContext context) {
    final iconTheme = IconTheme.of(context);
    final effectiveSize = size ?? iconTheme.size ?? 24.0;
    final effectiveColor = color ?? iconTheme.color;

    return SvgPicture.string(
      svg,
      width: effectiveSize,
      height: effectiveSize,
      colorFilter: effectiveColor != null ? ColorFilter.mode(effectiveColor, BlendMode.srcIn) : null,
    );
  }
}
