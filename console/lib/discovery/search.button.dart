import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library/api.dart' as api;

// SearchButton requests the peer-to-peer network locate media matching the
// current library search state - the free-text analog of KnownMediaLocator,
// which locates a specific already-catalogued item.
class SearchButton extends StatefulWidget {
  static const Widget _help = ds.Hint(Text("continuously scan the network for this search and feed matches into recommendations"));
  final media.MediaSearchState search;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;
  final Widget? label;
  final Widget help;

  const SearchButton({
    super.key,
    required this.search,
    this.locate = api.locate.create,
    this.label,
    this.help = _help,
  });

  @override
  State<SearchButton> createState() => _SearchButtonState();
}

class _SearchButtonState extends State<SearchButton> with ds.LoadingState {
  bool _queued = false;

  @override
  void initState() {
    super.initState();
    loading = false;
  }

  Future<void> _onPressed() async {
    setState(() {
      cause = ds.Error.zero;
    });

    final query = widget.search.next.query.trim();
    final mimetype = mimex.category(widget.search.next.mimetypes);
    if (query.isEmpty || mimetype.isEmpty) return;

    final options = [authn.request(authn.AuthzCache.meta(context))];

    await widget
        .locate(
          api.Locate.create()
            ..query = query
            ..mimetype = mimetype,
          options: options,
        )
        .then((v) {
          setState(() {
            _queued = true;
          });
        })
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final query = widget.search.next.query.trim();
    final mimetype = mimex.category(widget.search.next.mimetypes);
    if (query.isEmpty || mimetype.isEmpty) return ds.Empty;

    final label = widget.label ?? Text(_queued ? "monitoring" : "monitor");
    final icon = Icon(_queued ? Icons.radar_rounded : Icons.travel_explore_rounded);
    final btn = label == ds.Empty
        ? ds.LoadingIconButton(
            onPressed: _onPressed,
            icon: icon,
            tooltip: _queued ? "monitoring" : "monitor",
            disabled: _queued,
          )
        : ds.LoadingButton(
            Row(
              mainAxisSize: MainAxisSize.min,
              spacing: 8,
              children: [icon, label],
            ),
            onPressed: _onPressed,
            disabled: _queued,
          );

    return ds.Help(
      ds.Loading(
        cause: cause,
        btn,
      ),
      ds.Hint(Text("continuously scans the network for this search and feeds matches into recommendations")),
    );
  }
}
