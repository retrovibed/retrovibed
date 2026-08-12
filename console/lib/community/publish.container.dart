import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/media/media.known.pb.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/publish.content.dart';
import 'package:retrovibed/community/publish.metadata.dart';
import 'package:retrovibed/community/publish.confirmation.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class _StepIndicator extends StatelessWidget {
  final List<(Widget, Widget, Widget Function(_PublishContainerState))> steps;
  final int current;

  const _StepIndicator({required this.steps, required this.current});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children:
          steps.asMap().entries.map((entry) {
            final index = entry.key;
            final (_, indicator, _) = entry.value;
            final isActive = current == index;
            final isPast = current > index;

            return Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                if (index > 0)
                  Container(
                    width: 24,
                    height: 2,
                    color: isPast ? theme.colorScheme.primary : theme.colorScheme.outline,
                  ),
                Container(
                  padding: EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color:
                        isActive
                            ? theme.colorScheme.primary
                            : isPast
                            ? theme.colorScheme.primaryContainer
                            : Colors.transparent,
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: isActive || isPast ? theme.colorScheme.primary : theme.colorScheme.outline,
                    ),
                  ),
                  child: DefaultTextStyle(
                    style: theme.textTheme.labelSmall!.copyWith(
                      color:
                          isActive
                              ? theme.colorScheme.onPrimary
                              : isPast
                              ? theme.colorScheme.primary
                              : theme.colorScheme.outline,
                    ),
                    child: indicator,
                  ),
                ),
              ],
            );
          }).toList(),
    );
  }
}

class PublishContainer extends StatefulWidget {
  final VoidCallback onPublished;
  final VoidCallback onCancel;
  final Community? community;
  final media.FnMediaSearch search;
  final media.FnUploadRequest upload;

  const PublishContainer({
    super.key,
    required this.onPublished,
    required this.onCancel,
    this.community,
    this.search = media.media.search,
    this.upload = media.media.upload,
  });

  @override
  State<PublishContainer> createState() => _PublishContainerState();
}

class _PublishContainerState extends State<PublishContainer> with LoadingState {
  int _step = 0;
  Download? _download;
  Known? _known;
  Community? _community;
  Widget _cause = ds.Error.zero;

  @override
  void initState() {
    super.initState();
    _community = widget.community;
  }

  void _mutate(void Function(_PublishContainerState) fn) {
    setState(() {
      fn(this);
    });
  }

  void navigate(int delta) {
    _mutate((s) {
      s._step = (s._step + delta).clamp(0, _steps.length - 1);
    });
  }

  static final List<(Widget, Widget, Widget Function(_PublishContainerState))> _steps = [
    (
      const Text('Select Content'),
      const Text('Library'),
      (s) => PublishContent(
        search: s.widget.search,
        upload: s.widget.upload,
        onSelect:
            (item) => s._mutate((state) {
              state._download = item;
              state._step = 1;
            }),
      ),
    ),
    (
      const Text('Media Info'),
      const Text('Media'),
      (s) => PublishMetadata(
        download: s._download,
        onConfirm:
            (known) => s._mutate((state) {
              state._known = known;
              state._step = 2;
            }),
      ),
    ),
    (
      const Text('Confirm Publishing'),
      const Text('Confirm'),
      (s) => PublishConfirmation(
        download: s._download,
        community: s._community,
        knownMedia: s._known,
        onPublished: s.widget.onPublished,
      ),
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final (title, _, builder) = _steps[_step];
    final content = builder(this);

    return ds.Container(
      padding: defaults.padding,
      constraints: BoxConstraints(maxWidth: 700, maxHeight: 600),
      margin: defaults.padding / 2,
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              if (_step > 0) IconButton(onPressed: () => navigate(-1), icon: Icon(Icons.arrow_back)),
              DefaultTextStyle(
                style: theme.textTheme.headlineSmall!,
                textAlign: TextAlign.center,
                child: title,
              ),
              IconButton(onPressed: widget.onCancel, icon: Icon(Icons.close)),
            ],
          ),
          _StepIndicator(steps: _steps, current: _step),
          SizedBox(height: defaults.spacing * 2),
          Expanded(
            child: SingleChildScrollView(
              child: ds.ErrorScreen(content, cause: _cause),
            ),
          ),
        ],
      ),
    );
  }
}
