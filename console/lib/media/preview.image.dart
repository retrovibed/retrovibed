import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './media.pb.dart';
import './preview.unsupported.dart';

// the daemon already publishes a url for image mimetypes, so this costs no new transport.
class PreviewImage extends StatelessWidget {
  final Media current;

  const PreviewImage({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    return ds.Image.precache(
          context,
          current.image,
          headers: httpx.localheaders(current.image),
          fit: BoxFit.contain,
        ) ??
        PreviewUnsupported(current: current);
  }
}
