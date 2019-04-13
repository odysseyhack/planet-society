//
//  TransactionOverviewViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionVerificationStatus {
    case verified, unverified

    var color: UIColor {
        switch self {
        case .verified:
            return PHColors.topaz
        case .unverified:
            return PHColors.red
        }
    }
}

enum PHTableViewViewCellType {
    case notification(
        type: TransactionNotificationType,
        text: String)
    case warning(text: String)
    case description(date: Date, title: String, description: String)
    case transactionItem(item: TransactionItem)
    case form
}

protocol PHTableViewControllerDelegate: class {
    func didSelect(item: PHTableViewViewCellType)
    func didDeselect(item: PHTableViewViewCellType)
}

class PHTableViewController: UIViewController {

    // MARK: - Private properties

    lazy var tableView: UITableView = {

        let tableView = UITableView()
        tableView.translatesAutoresizingMaskIntoConstraints = false

        tableView.rowHeight = UITableView.automaticDimension
        tableView.estimatedRowHeight = 66
        tableView.separatorStyle = .none
        tableView.allowsMultipleSelection = true

        tableView.dataSource = self
        tableView.delegate = self

        tableView.register(
            TransactionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionTableViewCell.self))

        tableView.register(
            TransactionDescriptionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionDescriptionTableViewCell.self))

        tableView.register(
            TransactionNotificationTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionNotificationTableViewCell.self))

        tableView.register(
            FormTextInputCell.self,
            forCellReuseIdentifier: String(describing: FormTextInputCell.self))

        return tableView
    }()

    private var items: [PHTableViewViewCellType]

    // MARK: - Properties

    weak var delegate: PHTableViewControllerDelegate?

    // MARK: - Initialization

    init(
        title: String,
        items: [PHTableViewViewCellType]) {

        self.items = items

        super.init(nibName: nil, bundle: nil)

        self.title = title
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        // configure navigation bar
        navigationItem.title = title
        navigationController?.navigationBar.tintColor = PHColors.greyishBrown
        navigationItem.backBarButtonItem = UIBarButtonItem(title: "", style: .plain, target: nil, action: nil)

        view.addSubview(tableView)

        tableView.topAnchor.constraint(equalTo: view.topAnchor).isActive = true
        tableView.leftAnchor.constraint(equalTo: view.leftAnchor).isActive = true
        tableView.rightAnchor.constraint(equalTo: view.rightAnchor).isActive = true
        tableView.bottomAnchor.constraint(equalTo: view.bottomAnchor).isActive = true
    }
}

extension PHTableViewController: UITableViewDataSource {

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return items.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {

        switch items[indexPath.row] {
        case .notification(let type, let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionNotificationTableViewCell.self),
                for: indexPath) as! TransactionNotificationTableViewCell

            cell.configure(withType: type, andText: text)

            return cell

        case .description(let date, let title, let description):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionDescriptionTableViewCell.self),
                for: indexPath) as! TransactionDescriptionTableViewCell

            cell.configure(withDate: date, andTitle: title, andDescription: description)

            return cell

        case .transactionItem(let item):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: item.item,
                subtitle: (item.fields as NSArray).componentsJoined(by: ", "))
            cell.configure(withViewModel: viewModel)

            return cell

        case .warning(let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: text,
                subtitle: "")
            cell.configure(withViewModel: viewModel)

            return cell

        case .form:

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: FormTextInputCell.self),
                for: indexPath) as! FormTextInputCell

            cell.configure(withPlaceholder: "text") { text in
                print(text)
            }

            return cell
        }
    }
}

extension PHTableViewController: UITableViewDelegate {

    func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
        delegate?.didSelect(item: items[indexPath.row])

        switch items[indexPath.row] {
        case .transactionItem:
            tableView.cellForRow(at: indexPath)?.isSelected = true

        default:
            break
        }
    }

    func tableView(_ tableView: UITableView, didDeselectRowAt indexPath: IndexPath) {
        delegate?.didDeselect(item: items[indexPath.row])

        switch items[indexPath.row] {
        case .transactionItem:
            tableView.cellForRow(at: indexPath)?.isSelected = false

        default:
            break
        }
    }
}
