import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/mimex.dart' as mimex;
import './media.pb.dart';
import './preview.image.dart';
import './preview.pdf.dart';
import './preview.text.dart';

// renders content before it has been downloaded. playback already covers audio and video
// through PlayAction, so this dispatches on what is left.
abstract class Preview {
  // returns the action that opens the preview in ds.modals, or null when the mimetype has
  // no preview so callers can fall through to another action. trailing lets the caller add
  // actions beside close.
  static void Function()? modal(BuildContext context, Media current, {List<Widget> trailing = const []}) {
    final Widget body;

    if (mimex.isImage(current.mimetype)) {
      body = PreviewImage(current: current);
    } else if (current.mimetype == mimex.pdf) {
      body = PreviewPdf(current: current);
    } else if (mimex.isText(current.mimetype)) {
      body = PreviewText(current: current);
    } else {
      return null;
    }

    return () => ds.modals.push(
      context,
      ds.Card(
        constraints: ds.Defaults.modal(context),
        leading: [
          ds.Heading(
            Text(current.description),
            trailing: [
              ...trailing,
              IconButton(icon: const Icon(Icons.close), onPressed: () => ds.modals.push(context, null)),
            ],
          ),
        ],
        body,
      ),
    );
  }
}
