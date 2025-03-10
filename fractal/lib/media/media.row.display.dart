import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:media_kit/media_kit.dart' as mediakit;
import './api.dart' as api;
import './media.pb.dart';
import './player.dart';

void notimplemented(String s) {
  return print(s);
}

void rowtapdefault() => notimplemented("row tap not implemented");
void playtapdefault() => notimplemented("play tap not implemented");

class RowDisplay extends StatelessWidget {
  final Media media;
  final void Function() onTap;
  final void Function() onPlay;
  const RowDisplay({
    super.key,
    required this.media,
    this.onTap = rowtapdefault,
    this.onPlay = playtapdefault,
  });

  @override
  Widget build(BuildContext context) {
    const gap = SizedBox(width: 10.0);
    final theming = ds.Defaults.of(context);
    final vscreen = VideoScreen.of(context);

    return Container(
      padding: theming.padding,
      child: InkWell(
        onTap: onTap,
        child: Row(
          children: [
            const Icon(Icons.movie),
            gap,
            Expanded(
              child: Text(media.description, overflow: TextOverflow.ellipsis),
            ),
            gap,
            IconButton(
              icon: Icon(Icons.play_circle_outline_rounded),
              onPressed:
                  vscreen == null
                      ? null
                      : () {
                        vscreen.add(
                          mediakit.Media(api.media.download_uri(media.id)),
                        );
                      },
            ),
          ],
        ),
      ),
    );
  }
}
