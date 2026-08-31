import 'package:flutter/material.dart';
import 'package:retrovibed/mimex.dart' as mimex;
import './media.pb.dart';
import './preview.image.dart';
import './preview.pdf.dart';
import './preview.text.dart';
import './preview.unsupported.dart';

// renders content before it has been downloaded. playback already covers audio and video
// through PlayAction, so this dispatches on what is left.
class Preview extends StatelessWidget {
  final Media current;

  const Preview({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    if (mimex.isImage(current.mimetype)) {
      return PreviewImage(current: current);
    }

    if (current.mimetype == mimex.pdf) {
      return PreviewPdf(current: current);
    }

    if (mimex.isText(current.mimetype)) {
      return PreviewText(current: current);
    }

    return PreviewUnsupported(current: current);
  }
}
