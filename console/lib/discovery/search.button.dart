import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library/api.dart' as api;
import 'locate.p2p.prompt.dart' as p2p;

// SearchButton requests the peer-to-peer network locate media matching the
// current library search state - the free-text analog of KnownMediaLocator,
// which locates a specific already-catalogued item.
class SearchButton extends StatefulWidget {
  final media.MediaSearchState search;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;

  const SearchButton({
    super.key,
    required this.search,
    this.locate = api.locate.create,
  });

  @override
  State<SearchButton> createState() => _SearchButtonState();
}

class _SearchButtonState extends State<SearchButton> {
  bool _queued = false;
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> _onPressed() async {
    setState(() {
      _cause = ds.Error.zero;
    });

    final query = widget.search.next.query.trim();
    final mimetype = mimex.category(widget.search.next.mimetypes);
    if (query.isEmpty || mimetype.isEmpty) return;

    final options = [authn.request(authn.AuthzCache.meta(context))];

    try {
      final proceed = await p2p.ensureP2P(context, options: options);
      if (!proceed) return;

      await widget.locate(
        api.Locate.create()
          ..query = query
          ..mimetype = mimetype,
        options: options,
      );
      setState(() {
        _queued = true;
      });
    } catch (cause) {
      setState(() {
        _cause = ds.Error.unknown(cause, onTap: reseterr);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final query = widget.search.next.query.trim();
    final mimetype = mimex.category(widget.search.next.mimetypes);
    if (query.isEmpty || mimetype.isEmpty) return ds.Empty;

    return ds.Loading(
      cause: _cause,
      ds.LoadingButton(
        Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(_queued ? Icons.query_builder_rounded : Icons.travel_explore_rounded),
            const SizedBox(width: 8),
            Text(_queued ? "queued" : "discover"),
          ],
        ),
        onPressed: _onPressed,
        disabled: _queued,
      ),
    );
  }
}
