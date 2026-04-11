
import 'package:flutter/material.dart';

class Rating extends StatelessWidget {
  final double rating;
  final double size;
  final Color color;

  const Rating({
    Key? key,
    required this.rating,
    this.color = Colors.amber,
    this.size = 16,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ShaderMask(
      blendMode: BlendMode.srcATop, // Important for masking effect
      shaderCallback: (Rect bounds) {
        return LinearGradient(
          colors: [color, color],
        ).createShader(bounds);
      },
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            rating > 0.2 ? Icons.star : Icons.star_border,
            size: size,
            color: Colors.white, // Use white for the "mask"
          ),
          Icon(
            rating > 0.4 ? Icons.star : Icons.star_border,
            size: size,
            color: Colors.white, // Use white for the "mask"
          ),
          Icon(
            rating > 0.6 ? Icons.star : Icons.star_border,
            // index < rating ? Icons.star : Icons.star_border,
            size: size,
            color: Colors.white, // Use white for the "mask"
          ),
          Icon(
            rating > 0.8 ? Icons.star : Icons.star_border,
            // index < rating ? Icons.star : Icons.star_border,
            size: size,
            color: Colors.white, // Use white for the "mask"
          ),
          Icon(
            rating >= 1.0 ? Icons.star : Icons.star_border,
            // index < rating ? Icons.star : Icons.star_border,
            size: size,
            color: Colors.white, // Use white for the "mask"
          ),
        ]
      ),
    );
  }
}