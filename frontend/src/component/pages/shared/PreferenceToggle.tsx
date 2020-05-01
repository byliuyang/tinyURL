import React, { Component, ReactChild } from 'react';
import { Toggle } from '../../ui/Toggle';

interface IProps {
  toggleLabel: ReactChild;
  defaultIsEnabled: boolean;
  onToggleClick?: (enabled: boolean) => void;
}

const DEFAULT_PROPS = {
  defaultIsEnabled: false
};

export class PreferenceToggle extends Component<IProps> {
  static defaultProps: Partial<IProps> = DEFAULT_PROPS;

  render() {
    return (
      <div className={'preference-toggle'}>
        <Toggle
          defaultIsEnabled={this.props.defaultIsEnabled}
          onClick={this.props.onToggleClick}
        ></Toggle>
        <span className={'toggle-label'}>{this.props.toggleLabel}</span>
      </div>
    );
  }
}
