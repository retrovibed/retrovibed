import 'package:flutter/material.dart';
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/designkit.dart' as ds;
import './play.queue.dart';
import './playlist.dart' as internal;

class PlayerControlShuffle extends StatelessWidget {
  final ValueNotifier<PlayableMedia?> events;
  const PlayerControlShuffle(this.events, {Key? key}) : super(key: key);

  static IconData _icon(RangeFn mode) {
    if (mode == random) return Icons.shuffle;
    if (mode == acoustic) return Icons.graphic_eq;
    return Icons.queue;
  }

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<PlayableMedia?>(
      valueListenable: events,
      builder: (context, media, _) {
        switch (mimex.icon(media?.current.mimetype ?? mimex.binary)) {
          case mimex.icoaudio:
            final playlist = internal.Playlist.of(context);
            if (playlist == null) return ds.Empty;
            return ds.Help(
              IconButton(
                icon: Icon(_icon(playlist.autoqueue)),
                onPressed: () {
                  ds.modals
                      .asyncfn<RangeFn>(
                        context,
                        (completion) => _PlaybackStrategyMenu(
                          current: playlist.autoqueue,
                          onSelected: completion.complete,
                        ),
                      )
                      .then((mode) {
                        final current = media?.current;
                        if (current == null) return;
                        playlist.setPlaylist(
                          playlist.search.value.next,
                          current,
                          mode,
                          pos: playlist.player.state.position,
                        );
                      });
                },
              ),
              ds.Hint.multiline(const [
                Text("choose playback strategy"),
                Text("Search: play through the current search results, then fall back to random media"),
                Text("Random: skip the search results and play random media"),
                Text("Auto: play media that sounds similar to what's currently playing"),
              ]),
            );
          default:
            return ds.Empty;
        }
      },
    );
  }
}

class _PlaybackStrategyMenu extends StatelessWidget {
  static const _options = <(RangeFn, String)>[
    (search, 'Search'),
    (random, 'Random'),
    (acoustic, 'Auto'),
  ];

  final RangeFn current;
  final void Function(RangeFn mode) onSelected;
  const _PlaybackStrategyMenu({required this.current, required this.onSelected});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);
    return ds.Container(
      padding: defaults.padding,
      margin: defaults.margin,
      constraints: const BoxConstraints(maxWidth: 256),
      Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          for (final (mode, label) in _options)
            IconButton(
              tooltip: label,
              color: mode == current ? theme.colorScheme.primary : null,
              icon: Icon(PlayerControlShuffle._icon(mode)),
              onPressed: () => onSelected(mode),
            ),
        ],
      ),
    );
  }
}
