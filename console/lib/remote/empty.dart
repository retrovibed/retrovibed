import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/discovery.dart' as disc;
import 'api.dart' as remote;
import 'playlist.queue.dart';

class Empty extends StatelessWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final remote.RemoteControlSocket socket;
  final String sessionID;
  final remote.Sync current;
  final void Function(remote.Sync Function(remote.Sync c) mutated) onChange;
  final Future<void> Function()? Function(BuildContext, media.Media, media.MediaSearchResponse) onPlay;
  final lib.FnRecent apirecentlatest;
  final media.FnMediaSearch apisearch;

  const Empty(
    remote.Sync this.current, {
    Key? key,
    required this.search,
    required this.socket,
    required this.sessionID,
    required this.onPlay,
    required this.apisearch,
    this.onChange = ds.fnNoop,
    this.apirecentlatest = lib.recent.latest,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ValueListenableBuilder<media.MediaSearchState>(
      valueListenable: search,
      builder: (context, state, _) {
        final category = mimex.category(state.next.mimetypes);
        return PlaylistQueue(
          current,
          socket,
          key: const ValueKey("queue"),
          sessionId: sessionID,
          onChange: onChange,
          empty: ds.Loading(
            loading: current.token.isEmpty,
            maintainState: false,
            maintainAnimation: false,
            maintainSize: false,
            Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                ds.Heading(Text("Continue Watching", textAlign: TextAlign.center), padding: defaults.padding / 2),
                Expanded(
                  child: disc.RecentList(
                    margin: defaults.margin.copyWith(left: 0, right: 0),
                    padding: defaults.padding / 2,
                    latest: apirecentlatest,
                    host: current.library.hostname,
                    authz: httpx.Request.bearer(() => Future.value(current.token)),
                    onTap: (BuildContext context, media.RecentRecordRequest item) async {
                      return (onPlay(context, item.media, media.MediaSearchResponse(next: item.query)) ??
                          () async {})();
                    },
                    empty: lib.SearchMinimal(
                      key: const ValueKey("search-empty"),
                      empty: ds.Empty,
                      onPlay: onPlay,
                      apisearch: apisearch,
                      search: search,
                    ),
                    category,
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
