/**
 * 待办事项提醒接口定义
 */
export interface Todo {
  /**
   * 待办事项唯一标识符
   */
  id: string;
  
  /**
   * 待办事项标题
   */
  title: string;
  
  /**
   * 待办事项详情描述
   */
  detail: string;
  
  /**
   * 待办事项优先级
   * 可选值: 'low' | 'medium' | 'high'
   */
  priority: 'low' | 'medium' | 'high';
  
  /**
   * 是否启用提醒功能
   */
  remind: boolean;

  /**
   * 提醒时间
   * 如果不设置提醒，则为 null
   */
  remindTime: Date | null;
  
  /**
   * 待办事项创建时间
   */
  createdTime: Date;

  /**
   * 是否已完成
   */
  completed?: boolean;
}