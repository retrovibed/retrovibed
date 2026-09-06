import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'api.dart' as api;

/// What an installed publisher actually is: the identity the catalog and every
/// community selection hang off, the module it resolves to on disk, and the
/// fields an operator owns.
///
/// The id leads because it is the only thing about a publisher that is stable -
/// a clone arrives unnamed and shows its id until someone names it here.
/// actions are plain widgets so the card does not care that they happen to be
/// clone and delete.
class PublisherDetails extends StatefulWidget {
  final api.PluginPublisher current;
  final List<Widget> actions;
  final void Function(api.PluginPublisher updated) onChange;

  const PublisherDetails(
    this.current, {
    super.key,
    this.actions = const [],
    this.onChange = ds.fnNoop,
  });

  @override
  State<PublisherDetails> createState() => _PublisherDetails();
}

class _PublisherDetails extends State<PublisherDetails> with ds.LoadingState {
  // the endpoint takes the whole publisher, so edits accumulate here and are
  // saved as one row rather than a field at a time.
  late api.PluginPublisher _pending = widget.current.deepCopy();
  late api.PluginPublisher _saved = widget.current.deepCopy();

  Future<void> save() {
    // typing then tapping away, tapping away again, and submitting are all the
    // same edit; only an actual change is worth a round trip.
    if (_pending.description == _saved.description && _pending.mimetype == _saved.mimetype) {
      return Future.value();
    }

    setState(() => loading = true);
    return api.publishers
        .update(_pending, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _pending = v.publisher.deepCopy();
            _saved = v.publisher.deepCopy();
            cause = ds.Error.zero;
          });
          widget.onChange(v.publisher);
        })
        .catchError((cause) {
          setState(() => this.cause = ds.Error.unauthorized(cause, onTap: reseterr));
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((cause) {
          setState(() => this.cause = ds.Error.unknown(cause, onTap: reseterr));
        })
        .whenComplete(() => setState(() => loading = false));
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      padding: defaults.padding,
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        border: defaults.border,
        borderRadius: defaults.borderRadius,
      ),
      ds.Loading(
        loading: loading,
        cause: cause,
        Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          spacing: defaults.spacing,
          children: [
            Row(
              mainAxisSize: MainAxisSize.max,
              spacing: defaults.spacing,
              children: [
                Flexible(
                  child: SelectableText(
                    widget.current.id,
                    style: theme.textTheme.titleMedium,
                    maxLines: 1,
                  ),
                ),
                ...widget.actions,
              ],
            ),
            forms.Field(
              label: Text('Name'),
              input: TextFormField(
                initialValue: _pending.description,
                onChanged: (v) => _pending.description = v.trim(),
                onFieldSubmitted: (_) => save(),
                onTapOutside: (_) => save(),
                decoration: InputDecoration(
                  // an unnamed publisher displays as its id, which is what a
                  // clone looks like until it is named.
                  hintText: widget.current.id,
                  border: OutlineInputBorder(),
                ),
              ),
            ),
            forms.Field(
              label: Text('Mimetype'),
              input: inputs.Mimetype(
                value: _pending.mimetype,
                presets: [(label: Text('publish plugin') as Widget, value: mimex.publish)],
                onChanged: (v) {
                  _pending.mimetype = v;
                  save();
                },
              ),
              help: Text('Mimetype records what the plugin publishes; it is what a community matches against.'),
            ),
            forms.Field(
              label: Text('Path'),
              input: SelectableText(
                widget.current.path,
                style: theme.textTheme.bodySmall,
                maxLines: 1,
              ),
            ),
            ds.Timestamp.iso8601(
              widget.current.createdAt,
              leading: Text('Installed: ', style: theme.textTheme.bodySmall),
            ),
          ],
        ),
      ),
    );
  }
}
