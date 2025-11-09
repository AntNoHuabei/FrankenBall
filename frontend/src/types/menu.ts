import { WindowConfig } from "./window";
import {Component} from "vue";

export interface MenuItem {
  id: string;
  label: string;
  icon?: string | Component;
  type?: string;
  onClick?: (item: MenuItem) => void;
  windowConfig: WindowConfig;
  visible?: boolean; // 是否在菜单中显示，默认为true
}