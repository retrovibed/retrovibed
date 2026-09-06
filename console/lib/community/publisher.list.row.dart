import 'package:flutter/material.dart';
import 'package:path/path.dart' as p;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/ddisc/plugin/environment.editor.dart';
import 'api.dart' as api;
import './publisher.details.dart';

class PublisherRow extends StatelessWidget {
  final api.PluginPublisher current;
  final void Function(api.PluginPublisher deleted) onDelete;
  final void Function(api.PluginPublisher cloned) onClone;
  final void Function(api.PluginPublisher updated) onChange;

  const PublisherRow(
    this.current, {
    super.key,
    this.onDelete = ds.fnNoop,
    this.onClone = ds.fnNoop,
    this.onChange = ds.fnNoop,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.TableRow(
      key: ValueKey(current.id),
      expanded: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        spacing: defaults.spacing,
        children: [
          PublisherDetails(
            current,
            onChange: onChange,
            actions: [
              IconButton(
                icon: const Icon(Icons.copy_outlined),
                tooltip: 'clone this plugin into a second configuration',
                onPressed: () {
                  api.publishers
                      .clone(
                        current.id,
                        options: [authn.request(authn.AuthzCache.meta(context))],
                      )
                      .then((created) => onClone(created.publisher));
                },
              ),
              IconButton(
                icon: const Icon(Icons.delete_outline),
                onPressed: () {
                  final modal = ds.modals.of(context);
                  modal?.push(
                    ds.Confirmation.yesNo(
                      content: Text('Delete ${p.basename(current.path)}?'),
                      onCancel: (_) => modal.push(null),
                      onConfirm: (_) {
                        api.publishers
                            .delete(
                              current.id,
                              options: [authn.request(authn.AuthzCache.meta(context))],
                            )
                            .then((_) {
                              modal.push(null);
                              onDelete(current);
                            });
                      },
                    ),
                  );
                },
              ),
            ],
          ),
          // the same editor the search plugin list uses; the endpoint merges the
          // plugin's own declaration of the variables it understands over the
          // saved values, so a publisher nobody wrote a settings screen for
          // still renders a populated form.
          EnvironmentEditor.future(
            current.id,
            api.publisherenvironment.get(
              current.id,
              options: [authn.request(authn.AuthzCache.meta(context))],
            ),
            onChange: (content) {
              httpx
                  .withRetry(
                    () => api.publisherenvironment.update(
                      current.id,
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
          child: Text(
            current.description.isNotEmpty ? current.description : p.basename(current.path),
            overflow: TextOverflow.ellipsis,
            maxLines: 1,
          ),
        ),
        Text(current.mimetype),
      ],
    );
  }
}
