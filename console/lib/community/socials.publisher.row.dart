import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/ddisc/plugin/environment.editor.dart';
import 'api.dart' as api;
import './publisher.typography.dart';

// One publisher a community publishes through.
//
// The selection records a publisher id and nothing else, so the row resolves
// the plugin itself to have a name and a mimetype to show.
//
// Deliberately not PublisherRow: renaming, cloning and deleting are catalog
// wide, and here a delete would read as "stop using this" while actually
// uninstalling the plugin for every other community. Detaching is all this row
// offers, and it lives in the expanded region beside the configuration it
// applies to.
class SocialsPublisherRow extends StatefulWidget {
  final api.CommunityPublisher current;
  final api.FnPublishersFind find;
  final api.FnSocialsDisable disable;
  final void Function(api.CommunityPublisher detached) onDetach;

  const SocialsPublisherRow(
    this.current, {
    super.key,
    this.find = api.publishers.find,
    this.disable = api.socials.disable,
    this.onDetach = ds.fnNoop,
  });

  @override
  State<SocialsPublisherRow> createState() => _SocialsPublisherRow();
}

class _SocialsPublisherRow extends State<SocialsPublisherRow> {
  // started once and held: a row rebuilds whenever it expands or collapses,
  // and a future built in build() would refetch every time.
  Future<api.PluginPublisherFindResponse> _pending = Future.value(api.PluginPublisherFindResponse());
  Future<String> _environment = Future.value(EnvironmentEditor.zero);

  // the authz token comes from an inherited widget, which initState is too
  // early to read, hence the postframe.
  void _refresh() {
    final authz = [authn.request(authn.AuthzCache.meta(context))];
    setState(() {
      _pending = widget.find(widget.current.publisherId, options: authz);
      _environment = api.publisherenvironment.get(widget.current.publisherId, options: authz);
    });
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(_refresh);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.TableRow(
      key: ValueKey(widget.current.id),
      expanded: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        spacing: defaults.spacing,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              IconButton(
                icon: const Icon(Icons.link_off),
                tooltip: 'stop publishing this community through this plugin',
                onPressed: () {
                  widget
                      .disable(
                        widget.current.communityId,
                        widget.current.publisherId,
                        options: [authn.request(authn.AuthzCache.meta(context))],
                      )
                      .then((_) => widget.onDetach(widget.current));
                },
              ),
            ],
          ),
          // the plugin's configuration is per installation, not per community,
          // so this is the same editor and the same sidecar the catalog screen
          // writes - offered here because this is where it gets used.
          EnvironmentEditor.future(
            widget.current.publisherId,
            _environment,
            onChange: (content) {
              httpx
                  .withRetry(
                    () => api.publisherenvironment.update(
                      widget.current.publisherId,
                      content,
                      options: [authn.request(authn.AuthzCache.meta(context))],
                    ),
                  )
                  .catchError((cause) {
                    print("failed to update publisher environment ${cause}");
                    return content;
                  });
            },
          ),
        ],
      ),
      [
        Expanded(
          child: ds.future(api.PluginPublisherFindResponse(), _pending, (snapshot) {
            final publisher = snapshot.data?.publisher ?? api.PluginPublisher();
            return ds.ErrorScreen(
              // until it resolves the publisher is blank, and typography falls
              // back to the id - which is the only thing the row knows anyway.
              PublisherTypography(
                publisher.id.isEmpty ? (api.PluginPublisher()..id = widget.current.publisherId) : publisher,
                trailing: [Text(publisher.mimetype)],
              ),
              cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
            );
          }),
        ),
      ],
    );
  }
}
