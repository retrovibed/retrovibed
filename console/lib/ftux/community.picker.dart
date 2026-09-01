import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/community/community.pb.dart';
import 'api.dart';

/// Fetches the curated default community suggestions and lets the user pick
/// which ones to subscribe to. Calls [onDone] once the user has either
/// subscribed or skipped — the caller owns dismissing itself in response.
class CommunityPicker extends StatefulWidget {
  final VoidCallback onDone;
  const CommunityPicker({super.key, required this.onDone});

  @override
  State<CommunityPicker> createState() => _CommunityPickerState();
}

class _CommunityPickerState extends State<CommunityPicker> {
  bool _fetched = false;
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  List<Community> _suggestions = [];
  final Set<String> _selected = {};

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (_fetched) return;
    _fetched = true;

    final auth = [authn.request(authn.AuthzCache.meta(context))];
    ftux
        .defaults(options: auth)
        .then((resp) {
          setState(() {
            _suggestions = resp.community;
            _loading = false;
          });
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: () => setState(() => _cause = ds.Error.zero));
            _loading = false;
          });
        });
  }

  void _subscribe() {
    final auth = [authn.request(authn.AuthzCache.meta(context))];
    ftux.subscribe(_selected.toList(), options: auth).then((_) => widget.onDone());
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      padding: defaults.padding,
      margin: defaults.margin,
      constraints: const BoxConstraints(maxWidth: 512),
      ds.Loading(
        cause: _cause,
        loading: _loading,
        Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          spacing: defaults.spacing,
          children: [
            Text("Get Started", style: theme.textTheme.titleMedium),
            const Text(
              "Retrovibed can subscribe you to a few curated communities to get "
              "you started. Pick the ones you're interested in.",
            ),
            const Divider(),
            for (final suggestion in _suggestions)
              forms.Checkbox(
                Text(suggestion.description),
                description: Text(suggestion.url),
                value: _selected.contains(suggestion.id),
                onChanged: (v) => setState(() {
                  if (v ?? false) {
                    _selected.add(suggestion.id);
                  } else {
                    _selected.remove(suggestion.id);
                  }
                }),
              ),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              spacing: defaults.spacing,
              children: [
                TextButton(onPressed: widget.onDone, child: const Text('Skip')),
                ds.LoadingButton(const Text('Continue'), onPressed: () async => _subscribe()),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
